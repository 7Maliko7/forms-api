package structs

type SaveRequest struct {
	Id          uint32       `json:"id,omitempty"`
	Fields      []Field      `json:"fields,omitempty"`
	Attachments []Attachment `json:"attachment,omitempty"`
}

type SaveResponse struct {
	Uuid *string `json:"uuid,omitempty"`
}

type GetFormRequest struct {
	Uuid string `json:"uuid,omitempty"`
}

type GetFormResponse struct {
	Fields      []Field      `json:"fields,omitempty"`
	Attachments []Attachment `json:"attachment,omitempty"`
}

type GetFormListRequest struct {
	Limit  uint32 `json:"limit,omitempty"`
	Offset uint32 `json:"offset,omitempty"`
}

type GetFormListResponse struct {
	Forms []GetFormResponse `json:"forms,omitempty"`
}

type Field struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Data string `json:"data,omitempty"`
}

type Attachment struct {
	Uuid string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Data string `json:"data,omitempty"`
}
