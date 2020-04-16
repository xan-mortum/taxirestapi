package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/xan-mortum/taxirestapi/components"
	"github.com/xan-mortum/taxirestapi/gen/restapi/operations"
)

type TaxiHandler struct {
	Logger components.Logger
	DB     *components.DB
}

func NewTaxiHandler(logger components.Logger, db *components.DB) *TaxiHandler {
	return &TaxiHandler{
		Logger: logger,
		DB:     db,
	}
}

func (handler *TaxiHandler) RequestHandler(params operations.RequestParams) middleware.Responder {
	request := handler.DB.GetRequest()
	return operations.NewRequestOK().WithPayload(request)
}

func (handler *TaxiHandler) RequestsHandler(params operations.RequestsParams) middleware.Responder {
	requests := handler.DB.GetRequests()
	return operations.NewRequestsOK().WithPayload(requests)
}
