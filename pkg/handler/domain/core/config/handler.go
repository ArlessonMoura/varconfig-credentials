package config

import (
	"context"
	"net/http"
	dto "projeto-crud-credentials/dto/config"
	ports "projeto-crud-credentials/pkg/handler"

	"github.com/Wizzi-Cloud/restwrapper"
)

type Handler struct {
	svc ports.IVarConfigService
}

func NewHandler(svc ports.IVarConfigService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Handle(wrapper *restwrapper.Wrapper) {
	method := wrapper.RequestWrapper.Method()
	switch method {
	case http.MethodPost:
		h.Create(wrapper)
	case http.MethodGet:
		id, exists := wrapper.RequestWrapper.GetPathParam("id")
		if exists && id != nil && *id != "" {
			h.GetByID(wrapper)
		} else {
			h.List(wrapper)
		}
	case http.MethodPut:
		h.Update(wrapper)
	case http.MethodDelete:
		h.Delete(wrapper)
	default:
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) Create(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	var pathParams PathParams
	if err := wrapper.RequestWrapper.BindPathParams(&pathParams); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	var req dto.ConfigCreateRequestDTO
	if err := wrapper.RequestWrapper.BindBody(&req); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	// O Service recebe os IDs da URL + o Payload do Body
	result, err := h.svc.Create(ctx, pathParams.OrgID, pathParams.BenchmarkID, &req)
	if err != nil {
		errMsg := err.Error()
		if len(errMsg) > 17 && errMsg[:17] == "validation error:" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, errMsg)
			return
		}
		if errMsg == "name is required" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, errMsg)
			return
		}
		if errMsg == "configuration with this name already exists for this organization and benchmark" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusConflict, errMsg)
			return
		}
		wrapper.ResponseWrapper.WriteServerErrorResponse()
		return
	}

	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusCreated, result)
}

func (h *Handler) GetByID(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	var pathParams PathParams
	if err := wrapper.RequestWrapper.BindPathParams(&pathParams); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	// Validar ID adicional para garantir que não contém espaços
	if err := pathParams.Validate(); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.GetByID(ctx, pathParams.OrgID, pathParams.BenchmarkID, pathParams.ID)
	if err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusNotFound, "Configuração não encontrada")
		return
	}

	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusOK, result)
}

func (h *Handler) List(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	var pathParams PathParams
	if err := wrapper.RequestWrapper.BindPathParams(&pathParams); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.List(ctx, pathParams.OrgID, pathParams.BenchmarkID)
	if err != nil {
		wrapper.ResponseWrapper.WriteServerErrorResponse()
		return
	}

	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusOK, result)
}

func (h *Handler) Update(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	var pathParams PathParams
	if err := wrapper.RequestWrapper.BindPathParams(&pathParams); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	if pathParams.ID == "" {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, "id is required for update")
		return
	}

	var req dto.ConfigUpdateRequestDTO
	if err := wrapper.RequestWrapper.BindBody(&req); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	if err := pathParams.Validate(); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.Update(ctx, pathParams.OrgID, pathParams.BenchmarkID, pathParams.ID, &req)
	if err != nil {
		errMsg := err.Error()
		if len(errMsg) > 17 && errMsg[:17] == "validation error:" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, errMsg)
			return
		}
		if errMsg == "name is required" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, errMsg)
			return
		}
		if errMsg == "configuration with this name already exists for this organization and benchmark" {
			wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusConflict, errMsg)
			return
		}
		wrapper.ResponseWrapper.WriteServerErrorResponse()
		return
	}

	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusOK, result)
}

func (h *Handler) Delete(wrapper *restwrapper.Wrapper) {
	ctx := context.Background()

	var pathParams PathParams
	if err := wrapper.RequestWrapper.BindPathParams(&pathParams); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	if pathParams.ID == "" {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, "id is required for delete")
		return
	}

	if err := pathParams.Validate(); err != nil {
		wrapper.ResponseWrapper.WriteClientErrorResponse(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(ctx, pathParams.OrgID, pathParams.BenchmarkID, pathParams.ID); err != nil {
		wrapper.ResponseWrapper.WriteServerErrorResponse()
		return
	}

	wrapper.ResponseWrapper.WriteSuccessResponse(http.StatusNoContent, nil)
}
