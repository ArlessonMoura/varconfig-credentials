package benchmark

import (
	"context"
	"net/http"
	ports "projeto-crud-credentials/pkg/handler"

	"github.com/Wizzi-Cloud/restwrapper"
)

type Handler struct {
	svc ports.IBenchmarkSchemaService
}

func NewHandler(svc ports.IBenchmarkSchemaService) *Handler {
	return &Handler{
		svc: svc,
	}
}

// Handle implementa a interface IRestHandler do restwrapper
func (h *Handler) Handle(wrapper *restwrapper.Wrapper) {
	method := wrapper.RequestWrapper.Method()
	switch method {
	case http.MethodGet:
		benchmarkId, exists := wrapper.RequestWrapper.GetPathParam("benchmarkId")
		if exists && benchmarkId != nil && *benchmarkId != "" {
			h.GetByID(wrapper)
		} else {
			h.List(wrapper)
		}
	default:
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) GetByID(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	// Binding dos path parameters necessários para este endpoint
	var pathParams PathParams
	if err := wrapper.RequestWrapper.BindPathParams(&pathParams); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	schema, err := h.svc.GetByID(ctx, pathParams.BenchmarkID)
	if err != nil {
		if err.Error() == "schema not found" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusNotFound, "Schema not found")
			return
		}
		wrapper.ResponseWrapper.WriteServerErrorResponse()
		return
	}
	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusOK, schema)
}

func (h *Handler) List(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	// Para List não precisamos de path parameters obrigatórios
	// Chamamos diretamente o service
	schemas, err := h.svc.List(ctx)
	if err != nil {
		wrapper.ResponseWrapper.WriteServerErrorResponse()
		return
	}

	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusOK, schemas)
}
