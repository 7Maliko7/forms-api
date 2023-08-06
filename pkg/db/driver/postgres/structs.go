package postgres

type DbField struct {
	Uuid string `db:"data_uuid"`
	Name string `db:"field_name"`
	Type string `db:"field_type"`
	Data string `db:"field_data"`
}

type DbAttachment struct {
	Uuid string `db:"file_uuid"`
	Name string `db:"file_name"`
	Type string `db:"file_type"`
}

type DbForm struct {
	Uuid       string `db:"form_uuid"`
	Id         string `db:"form_id"`
	CreatedAt  string `db:"created_at"`
	Fields     []DbField
	Attachment []DbAttachment
}
