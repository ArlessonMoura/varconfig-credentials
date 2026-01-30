package varconfig

import (
	"net/http"
	dto "projeto-crud-credencials/dto/varconfig_dto"
	service "projeto-crud-credencials/pkg/handler"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.IVarConfigService
}

func NewHandler(service service.IVarConfigService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) List(c *gin.Context) {
	orgID := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")

	// Chamada ao service usando o contrato de DTO de resposta
	result, err := h.service.List(c.Request.Context(), orgID, benchmarkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Create(c *gin.Context) {
	orgID := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")

	var req dto.CreateVarConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	// Validação adicional do payload
	if err := ValidateCreateAndUpdateRequest(&PathParameter{Payload: req.Payload}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// O Service recebe os IDs da URL + o Payload do Body
	result, err := h.service.Create(c.Request.Context(), orgID, benchmarkID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetByID(c *gin.Context) {
	orgID := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")
	id := c.Param("id")

	result, err := h.service.GetByID(c.Request.Context(), orgID, benchmarkID, id)
	if err != nil {
		// Em produção, aqui usaríamos o mapeamento de erros do internal/common
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuração não encontrada"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Update(c *gin.Context) {
	orgID := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")
	id := c.Param("id")

	var req dto.UpdateVarConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validação adicional do payload
	if err := ValidateCreateAndUpdateRequest(&PathParameter{Payload: req.Payload}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Update(c.Request.Context(), orgID, benchmarkID, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Delete(c *gin.Context) {
	orgID := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), orgID, benchmarkID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
