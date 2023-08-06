package forms

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/docker/distribution/uuid"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	svc "github.com/7Maliko7/forms-api/internal/service"
	"github.com/7Maliko7/forms-api/internal/transport/structs"
	"github.com/7Maliko7/forms-api/pkg/broker"
	"github.com/7Maliko7/forms-api/pkg/db"
	statuses "github.com/7Maliko7/forms-api/pkg/errors"
	"github.com/7Maliko7/forms-api/pkg/storage"
)

const (
	logKeyMethod = "method"
	logKeyErr    = "err"
)

type Service struct {
	Repository db.Databaser
	Logger     log.Logger
	Broker     broker.Broker
	Storage    storage.Storager
}

type brokerMessage struct {
	Uuid   uuid.UUID       `json:"uuid"`
	Id     uint32          `json:"id"`
	Fields []structs.Field `json:"fields"`
}

func NewService(rep db.Databaser, logger log.Logger, b broker.Broker, s storage.Storager) svc.Service {
	return &Service{
		Repository: rep,
		Logger:     logger,
		Broker:     b,
		Storage:    s,
	}
}

func (s *Service) Save(ctx context.Context, id uint32, fields []structs.Field, attachment []structs.Attachment) (*uuid.UUID, error) {
	logger := log.With(s.Logger, logKeyMethod, "Create")

	dbFields := make([]db.Field, 0, len(fields))
	for _, v := range fields {
		dbFields = append(dbFields, db.Field{
			Name: v.Name,
			Type: v.Type,
			Data: v.Data,
		})
	}

	formUuid, err := s.Repository.Save(ctx, id, dbFields)
	if err != nil {
		level.Error(logger).Log("repository", err.Error())
		return nil, statuses.FailedRequest
	}
	level.Debug(logger).Log("New formUuid", formUuid.String())

	for _, at := range attachment {
		fileUuid := uuid.Generate()
		_, err = s.Repository.SaveAttachment(fileUuid, formUuid, at.Name, at.Type)
		if err != nil {
			level.Error(logger).Log("repository", err.Error())
			return nil, statuses.FailedRequest
		}
		data, err := base64.StdEncoding.DecodeString(at.Data)
		if err != nil {
			level.Error(logger).Log("repository", err.Error())
			return nil, statuses.FailedRequest
		}
		err = s.Storage.SaveFile(ctx, fileUuid.String(), data)
		if err != nil {
			level.Error(logger).Log("repository", err.Error())
			return nil, statuses.FailedRequest
		}
	}

	b := brokerMessage{
		Uuid:   formUuid,
		Id:     id,
		Fields: fields,
	}
	body, err := json.Marshal(b)
	if err != nil {
		level.Error(logger).Log("broker", err.Error())
		return nil, statuses.FailedRequest
	}
	err = s.Broker.Publish(ctx, body, "forms.data", "")
	if err != nil {
		level.Error(logger).Log("broker", err.Error())
		return nil, statuses.FailedRequest
	}

	return &formUuid, nil
}

func (s *Service) GetForm(ctx context.Context, form uuid.UUID) (*structs.GetFormResponse, error) {
	logger := log.With(s.Logger, logKeyMethod, "Get form")
	Form, err := s.Repository.GetForm(form)
	if err != nil {
		level.Error(logger).Log("repository", err.Error())
		return nil, statuses.FailedRequest
	}
	level.Debug(logger).Log("New Form", Form)

	result := &structs.GetFormResponse{}
	result.Fields = make([]structs.Field, 0, len(Form.Fields))
	for _, v := range Form.Fields {
		result.Fields = append(result.Fields, structs.Field{Name: v.Name, Type: v.Type, Data: v.Data})
	}
	result.Attachments = make([]structs.Attachment, 0, len(Form.Attachment))
	for _, at := range Form.Attachment {
		result.Attachments = append(result.Attachments, structs.Attachment{Name: at.Name, Type: at.Type, Uuid: at.Uuid.String()})
	}
	return result, nil

}

func (s *Service) GetFormList(ctx context.Context, limit, offset uint32) (*structs.GetFormListResponse, error) {
	logger := log.With(s.Logger, logKeyMethod, "Get formList")
	Form, err := s.Repository.GetFormList(limit, offset)
	if err != nil {
		level.Error(logger).Log("repository", err.Error())
		return nil, statuses.FailedRequest
	}
	level.Debug(logger).Log("New FormList", Form)

	result := &structs.GetFormListResponse{}
	result.Forms = make([]structs.GetFormResponse, 0, len(Form))
	for _, v := range Form {
		form := structs.GetFormResponse{Fields: make([]structs.Field, 0, len(v.Fields))}
		for _, field := range v.Fields {
			form.Fields = append(form.Fields, structs.Field{Name: field.Name, Type: field.Type, Data: field.Data})
		}
		attachment := structs.GetFormResponse{Attachments: make([]structs.Attachment, 0, len(v.Attachment))}
		for _, at := range v.Attachment {
			attachment.Attachments = append(attachment.Attachments, structs.Attachment{Name: at.Name, Type: at.Type, Uuid: at.Uuid.String()})
		}
		result.Forms = append(result.Forms, form, attachment)

	}
	return result, nil

}
