package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/7Maliko7/forms-api/internal/transport"
	"github.com/7Maliko7/forms-api/internal/transport/structs"
	statuses "github.com/7Maliko7/forms-api/pkg/errors"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

// NewService wires Go kit endpoints to the HTTP transport.
func NewService(
	svcEndpoints transport.Endpoints, options []kithttp.ServerOption, logger log.Logger,
) http.Handler {
	var (
		r            = mux.NewRouter()
		errorLogger  = kithttp.ServerErrorLogger(logger)
		errorEncoder = kithttp.ServerErrorEncoder(encodeErrorResponse)
	)
	options = append(options, errorLogger, errorEncoder)
	r.Methods("POST").Path("/save").Handler(kithttp.NewServer(
		svcEndpoints.Save,
		decodeSaveRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/get").Handler(kithttp.NewServer(
		svcEndpoints.GetForm,
		decodeGetFormRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/list").Handler(kithttp.NewServer(
		svcEndpoints.GetFormList,
		decodeGetFormListRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeSaveRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req structs.SaveRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, statuses.InvalidRequest
	}

	return req, nil
}

func decodeGetFormRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req structs.GetFormRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, statuses.InvalidRequest
	}

	return req, nil
}

func decodeGetFormListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req structs.GetFormListRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, statuses.InvalidRequest
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeErrorResponse(ctx, e.error(), w)

		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case statuses.InvalidRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
