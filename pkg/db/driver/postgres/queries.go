package postgres

const (
	AddFormQuery       = `/*NO LOAD BALANCE*/ select * from operations.add_form($1::smallint) as form_uuid;`
	AddDataQuery       = `/*NO LOAD BALANCE*/ select * from operations.add_data($1::uuid,$2,$3,$4) as data_uuid;`
	AddAttachmentQuery = `/*NO LOAD BALANCE*/ select * from operations.add_attachment($1::uuid,$2::uuid,$3,$4) as file_uuid;`
	GetFormQuery       = `select * from operations.get_form($1::uuid);`
	GetAttachmentQuery = `select * from operations.get_attachment($1::uuid);`
	GetFormListQuery   = `select * from operations.get_form_list($1,$2);`
)
