package varconfig

// Importar DTOs do pacote dto
import (
	"net/http"
	"strconv"

	varconfigdto "projeto-crud-credencials/dto"
	svcvarconfig "projeto-crud-credencials/internal/service/domain/varconfig"

	"github.com/gin-gonic/gin"
)

// Handler estrutura para manipular requisições de VarConfig
type Handler struct {
	service VarConfigUseCase
}

// NewHandler cria uma nova instância do handler
func NewHandler(service VarConfigUseCase) *Handler {
	return &Handler{
		service: service,
	}
}

// ListByBenchmark (GET) lista todos os VarConfigs de um benchmark
func (h *Handler) ListByBenchmark(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgId inválido"})
		return
	}

	configs, err := h.service.ListByBenchmark(c.Request.Context(), orgID, benchmarkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": configs})
}

// Create (POST) cria um novo VarConfig
func (h *Handler) Create(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgId inválido"})
		return
	}

	var req varconfigdto.CreateVarConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateCreateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := svcvarconfig.VarConfig{
		OrgID:       orgID,
		BenchmarkID: benchmarkID,
		Payload:     req.Payload,
	}

	result, err := h.service.Create(c.Request.Context(), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetByID (GET) obtém um VarConfig específico
func (h *Handler) GetByID(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")
	idStr := c.Param("id")

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgId inválido"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	config, err := h.service.GetByID(c.Request.Context(), orgID, benchmarkID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// Update (PUT/PATCH) atualiza um VarConfig
func (h *Handler) Update(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")
	idStr := c.Param("id")

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgId inválido"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req varconfigdto.UpdateVarConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateUpdateRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := svcvarconfig.VarConfig{
		ID:          id,
		OrgID:       orgID,
		BenchmarkID: benchmarkID,
		Payload:     req.Payload,
	}

	result, err := h.service.Update(c.Request.Context(), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Delete (DELETE) remove um VarConfig
func (h *Handler) Delete(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	benchmarkID := c.Param("benchmark_id")
	idStr := c.Param("id")

	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "orgId inválido"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err = h.service.Delete(c.Request.Context(), orgID, benchmarkID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
