// Package benchmark_schema provides HTTP handlers for managing benchmark schemas.
package benchmark_schema

import (
	"errors"
	"net/http"
	service "projeto-crud-credencials/pkg/handler"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	retrievalService service.IBenchmarkSchemaService
}

func NewHandler(retrievalService service.IBenchmarkSchemaService) *Handler {
	return &Handler{
		retrievalService: retrievalService,
	}
}

// GetByID retorna um schema específico pelo benchmarkId
func (h *Handler) GetByID(c *gin.Context) {
	benchmarkID, isValid := ValidateBenchmarkID(c)
	if !isValid {
		return
	}

	schema, err := h.retrievalService.GetSchemaByID(c.Request.Context(), benchmarkID)
	if err != nil {
		if errors.Is(err, errors.New("schema not found")) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Schema not found",
				"details": "No benchmark schema found with the provided ID",
			})
			return
		}

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
	schemas, err := h.retrievalService.ListAllSchemas(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas)
}
