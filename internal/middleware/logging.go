package middleware

import (
	"context"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/go-kit/kit/log"

	"github.com/7Maliko7/forms-api/internal/service"
	"github.com/7Maliko7/forms-api/internal/transport/structs"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next service.Service) service.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   service.Service
	logger log.Logger
}

func (mw loggingMiddleware) Save(ctx context.Context, id uint32, fields []structs.Field, attachment []structs.Attachment) (*uuid.UUID, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Save", "formId", id, "fields", fields, "duration", time.Since(begin), "err")
	}(time.Now())
	return mw.next.Save(ctx, id, fields, attachment)
}

func (mw loggingMiddleware) GetForm(ctx context.Context, form uuid.UUID) (*structs.GetFormResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetForm", "duration", time.Since(begin), "err")
	}(time.Now())
	return mw.next.GetForm(ctx, form)
}

func (mw loggingMiddleware) GetFormList(ctx context.Context, limit, offset uint32) (*structs.GetFormListResponse, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetFormList", "duration", time.Since(begin), "err")
	}(time.Now())
	return mw.next.GetFormList(ctx, limit, offset)
}
