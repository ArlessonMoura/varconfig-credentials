package routes

import (
	"github.com/Wizzi-Cloud/restwrapper/handler"
	"github.com/gin-gonic/gin"
)

func SetupBenchmarkSchemaRoutes(router *gin.Engine, h *handler.Handler) {
	schemaGroup := router.Group("/compliance/benchmark/schema")
	{
		schemaGroup.GET("", h.HandleGin)

		schemaGroup.GET("/:benchmarkId", h.HandleGin)
	}
}
