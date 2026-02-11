package routes

import (
	benchmarkSchemaHandler "projeto-crud-credencials/pkg/handler/domain/compliance/benchmark_schema"
	"projeto-crud-credencials/pkg/handler/domain/core/varconfig"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *varconfig.Handler, bsHandler *benchmarkSchemaHandler.Handler) *gin.Engine {
	router := gin.Default()
	SetupVarConfigRoutes(router, h)
	SetupBenchmarkSchemaRoutes(router, bsHandler)

	return router
}
