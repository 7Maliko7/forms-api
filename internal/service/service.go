package service

import (
	"context"

	"github.com/7Maliko7/forms-api/internal/transport/structs"
	"github.com/docker/distribution/uuid"
)

type Service interface {
	Save(ctx context.Context, id uint32, fields []structs.Field, attachment []structs.Attachment) (*uuid.UUID, error)
	GetForm(ctx context.Context, form uuid.UUID) (*structs.GetFormResponse, error)
	GetFormList(ctx context.Context, limit, offset uint32) (*structs.GetFormListResponse, error)
}
