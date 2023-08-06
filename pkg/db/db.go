package db

import (
	"context"

	"github.com/docker/distribution/uuid"
)

type Databaser interface {
	Save(ctx context.Context, id uint32, fields []Field) (uuid.UUID, error)
	SaveAttachment(file uuid.UUID, form uuid.UUID, name string, fileType string) (uuid.UUID, error)
	GetForm(form uuid.UUID) (Form, error)
	GetFormList(limit, offset uint32) ([]Form, error)
}

type Field struct {
	Name string
	Type string
	Data string
}

type Attachment struct {
	Uuid uuid.UUID
	Name string
	Type string
}

type Form struct {
	Uuid   uuid.UUID
	Fields []Field
	Attachment []Attachment
}
