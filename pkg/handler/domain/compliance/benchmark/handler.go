package benchmark

import (
	"context"
	"net/http"
	dto "projeto-crud-credentials/dto/benchmark"
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
	case http.MethodPost:
		h.Create(wrapper)
	case http.MethodGet:
		benchmarkId, exists := wrapper.RequestWrapper.GetPathParam("benchmarkId")
		if exists && benchmarkId != nil && *benchmarkId != "" {
			h.GetByID(wrapper)
		} else {
			h.List(wrapper)
		}
	case http.MethodDelete:
		h.Delete(wrapper)
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

func (h *Handler) Create(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	var req dto.BenchmarkCreateRequestDTO
	if err := wrapper.RequestWrapper.BindBody(&req); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	result, err := h.svc.Create(ctx, &req)
	if err != nil {
		if err.Error() == "name is required" || err.Error() == "version is required" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "schema with this name already exists" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusConflict, err.Error())
			return
		}
		wrapper.ResponseWrapper.WriteServerErrorResponse()
		return
	}
	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusCreated, result)
}

func (h *Handler) Delete(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	// Binding dos path parameters necessários para este endpoint
	var pathParams PathParams
	if err := wrapper.RequestWrapper.BindPathParams(&pathParams); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	// Validar ID
	if err := pathParams.Validate(); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(ctx, pathParams.BenchmarkID); err != nil {
		if err.Error() == "schema not found" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusNotFound, "Schema not found")
			return
		}
		wrapper.ResponseWrapper.WriteServerErrorResponse()
		return
	}

	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusNoContent, nil)
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
