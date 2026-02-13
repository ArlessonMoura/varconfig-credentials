package benchmark_schema

import (
	"net/http"
	ports "projeto-crud-credencials/pkg/handler"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc ports.IBenchmarkSchemaService
}

func NewHandler(svc ports.IBenchmarkSchemaService) *Handler {
	return &Handler{
		svc: svc,
	}
}

// GetByID retorna um schema específico pelo benchmarkId
func (h *Handler) GetByID(c *gin.Context) {
	benchmarkID, isValid := ValidateBenchmarkID(c)
	if !isValid {
		return
	}

	schema, err := h.svc.GetByID(c.Request.Context(), benchmarkID)
	if err != nil {
		// if errors.Is(err, errors.New("schema not found")) {
		// 	c.JSON(http.StatusNotFound, gin.H{
		// 		"error":   "Schema not found",
		// 		"details": "No benchmark schema found with the provided ID",
		// 	})
		// 	return
		// }

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schema)
}

// List retorna todos os schemas disponíveis
func (h *Handler) List(c *gin.Context) {
	schemas, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas)
}
