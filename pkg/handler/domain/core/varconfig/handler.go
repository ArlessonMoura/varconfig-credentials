package varconfig

import (
	"net/http"
	dto "projeto-crud-credencials/dto/varconfig"
	ports "projeto-crud-credencials/pkg/handler"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ports.IVarConfigService
}

func NewHandler(svc ports.IVarConfigService) *Handler {
	return &Handler{
		svc: svc,
	}
}


func (h *Handler) Create(c *gin.Context) {
	orgID := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")

	var req dto.VarConfigCreationPayload
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
	result, err := h.svc.Create(c.Request.Context(), orgID, benchmarkID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}


func (h *Handler) List(c *gin.Context) {
	orgID := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")

	// Validação dos parâmetros de rota
	if err := ValidatePathParams(orgID, benchmarkID, ""); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chamada ao svc usando o contrato de DTO de resposta
	result, err := h.svc.List(c.Request.Context(), orgID, benchmarkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}



func (h *Handler) GetByID(c *gin.Context) {
	orgID := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")
	id := c.Param("id")

	// Validação dos parâmetros de rota
	if err := ValidatePathParams(orgID, benchmarkID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.GetByID(c.Request.Context(), orgID, benchmarkID, id)
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

	var req dto.VarConfigUpdatePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validação adicional do payload
	if err := ValidateCreateAndUpdateRequest(&PathParameter{Payload: req.Payload}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Update(c.Request.Context(), orgID, benchmarkID, id, req)
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

	// Validação dos parâmetros de rota
	if err := ValidatePathParams(orgID, benchmarkID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), orgID, benchmarkID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
