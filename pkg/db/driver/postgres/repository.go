package postgres

import (
	"context"
	"database/sql"

	_ "github.com/cockroachdb/cockroach-go/crdb"
	"github.com/docker/distribution/uuid"

	"github.com/7Maliko7/forms-api/pkg/db"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

func (repo *Repository) Save(ctx context.Context, id uint32, fields []db.Field) (uuid.UUID, error) {
	var formUuid string
	err := repo.db.QueryRow(AddFormQuery, id).Scan(&formUuid)
	if err != nil {
		return uuid.UUID{}, err
	}

	for _, v := range fields {
		var dataUuid string
		err := repo.db.QueryRow(AddDataQuery, formUuid, v.Name, v.Type, v.Data).Scan(&dataUuid)
		if err != nil {
			return uuid.UUID{}, err
		}
	}

	saveUuid, err := uuid.Parse(formUuid)
	if err != nil {
		return uuid.UUID{}, err
	}

	return saveUuid, nil
}

func (repo *Repository) SaveAttachment(file uuid.UUID, form uuid.UUID, name string, fileType string) (uuid.UUID, error) {
	var fileUuid string
	err := repo.db.QueryRow(AddAttachmentQuery, file.String(), form.String(), name, fileType).Scan(&fileUuid)
	if err != nil {
		return uuid.UUID{}, err
	}
	fUuid, err := uuid.Parse(fileUuid)
	if err != nil {
		return uuid.UUID{}, err
	}
	return fUuid, nil
}

func (repo *Repository) GetForm(form uuid.UUID) (db.Form, error) {
	var (
		result       db.Form
		dbField      []DbField
		dbAttachment []DbAttachment
	)

	rows, err := repo.db.Query(GetFormQuery, form.String())
	if err != nil {
		return db.Form{}, err
	}
	for rows.Next() {
		var fl DbField
		err = rows.Scan(&fl.Uuid, &fl.Name, &fl.Type, &fl.Data)
		if err != nil {
			return db.Form{}, err
		}

		dbField = append(dbField, fl)
	}
	result.Fields = make([]db.Field, 0, len(dbField))
	for _, v := range dbField {
		result.Fields = append(result.Fields, db.Field{Name: v.Name, Type: v.Type, Data: v.Data})
	}

	rowsAt, err := repo.db.Query(GetAttachmentQuery, form.String())
	if err != nil {
		return db.Form{}, err
	}
	for rowsAt.Next() {
		var at DbAttachment
		err = rowsAt.Scan(&at.Uuid, &at.Name, &at.Type)
		if err != nil {
			return db.Form{}, err
		}

		dbAttachment = append(dbAttachment, at)
	}
	result.Attachment = make([]db.Attachment, 0, len(dbAttachment))
	for _, f := range dbAttachment {
		uuidAttachment, err := uuid.Parse(f.Uuid)
		if err != nil {
			return db.Form{}, err
		}
		result.Attachment = append(result.Attachment, db.Attachment{Uuid: uuidAttachment, Name: f.Name, Type: f.Type})
	}

	return result, nil
}

func (repo *Repository) GetFormList(limit, offset uint32) ([]db.Form, error) {
	var dbForm []DbForm
	rows, err := repo.db.Query(GetFormListQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	dbForm = make([]DbForm, 0, 2)
	for rows.Next() {
		var f DbForm
		err = rows.Scan(&f.Uuid, &f.Id, &f.CreatedAt)
		if err != nil {
			return nil, err
		}

		dbForm = append(dbForm, f)
	}

	result := make([]db.Form, 0, len(dbForm))

	for _, v := range dbForm {
		rows, err := repo.db.Query(GetFormQuery, v.Uuid)
		if err != nil {
			return nil, err
		}

		var dbFields []DbField
		dbFields = make([]DbField, 0, 2)
		for rows.Next() {
			var fl DbField

			err = rows.Scan(&fl.Uuid, &fl.Name, &fl.Type, &fl.Data)
			if err != nil {
				return nil, err
			}

			dbFields = append(dbFields, fl)
		}

		uuidForm, err := uuid.Parse(v.Uuid)
		if err != nil {
			return nil, err
		}
		resultFields := make([]db.Field, 0, len(dbFields))
		for _, f := range dbFields {
			resultFields = append(resultFields, db.Field{Name: f.Name, Type: f.Type, Data: f.Data})
		}

		var dbAttachment []DbAttachment
		rowsAt, err := repo.db.Query(GetAttachmentQuery, v.Uuid)
		if err != nil {
			return nil, err
		}
		dbAttachment = make([]DbAttachment, 0, 2)
		for rowsAt.Next() {
			var at DbAttachment
			err = rowsAt.Scan(&at.Uuid, &at.Name, &at.Type)
			if err != nil {
				return nil, err
			}

			dbAttachment = append(dbAttachment, at)
		}

		resultAttachment := make([]db.Attachment, 0, len(dbAttachment))
		for _, a := range dbAttachment {
			uuidAttachment, err := uuid.Parse(a.Uuid)
			if err != nil {
				return nil, err
			}
			resultAttachment = append(resultAttachment, db.Attachment{Uuid: uuidAttachment, Name: a.Name, Type: a.Type})
		}

		result = append(result, db.Form{
			Uuid:       uuidForm,
			Fields:     resultFields,
			Attachment: resultAttachment,
		})

	}

	return result, nil
}
