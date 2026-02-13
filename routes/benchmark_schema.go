package routes

import (
	"projeto-crud-credencials/pkg/handler/domain/compliance/benchmark_schema"

	"github.com/gin-gonic/gin"
)

// SetupBenchmarkSchemaRoutes registra as rotas de consulta de schemas de benchmark
func SetupBenchmarkSchemaRoutes(router *gin.Engine, h *benchmark_schema.Handler) {
	schemaGroup := router.Group("/compliance/benchmark/schema")
	{
		schemaGroup.GET("", h.List)

		schemaGroup.GET("/:benchmarkId", h.GetByID)
	}
}
