package transport

import (
	"context"

	"github.com/docker/distribution/uuid"
	"github.com/go-kit/kit/endpoint"

	"github.com/7Maliko7/forms-api/internal/service"
	"github.com/7Maliko7/forms-api/internal/transport/structs"
)

type Endpoints struct {
	Save        endpoint.Endpoint
	GetForm     endpoint.Endpoint
	GetFormList endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Save:        makeSaveEndpoint(s),
		GetForm:     makeGetFormEndpoint(s),
		GetFormList: makeGetFormListEndpoint(s),
	}
}

func makeSaveEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(structs.SaveRequest)

		uuid, err := s.Save(ctx, req.Id, req.Fields, req.Attachments)
		if err != nil {
			return structs.SaveResponse{Uuid: nil}, err
		}

		var result string
		if uuid != nil {
			result = uuid.String()
		}
		return structs.SaveResponse{Uuid: &result}, nil
	}
}

func makeGetFormEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(structs.GetFormRequest)
		uu, err := uuid.Parse(req.Uuid)
		if err != nil {
			return structs.GetFormResponse{}, err
		}
		form, err := s.GetForm(ctx, uu)
		if err != nil {
			return structs.GetFormResponse{}, err
		}
		return form, nil
	}
}

func makeGetFormListEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(structs.GetFormListRequest)
		formList, err := s.GetFormList(ctx, req.Limit, req.Offset)
		if err != nil {
			return structs.GetFormListResponse{}, err
		}
		return formList, nil
	}
}
