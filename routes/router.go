package routes

import (
	"projeto-crud-credentials/pkg/handler/domain/compliance/benchmark"
	"projeto-crud-credentials/pkg/handler/domain/core/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(hConfig *config.Handler, hBench *benchmark.Handler) *gin.Engine {
	router := gin.Default()
	SetupVarConfigRoutes(router, hConfig)
	SetupBenchmarkSchemaRoutes(router, hBench)

	return router
}
