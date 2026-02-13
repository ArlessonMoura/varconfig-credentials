package routes

import (
	benchmarkSchema "projeto-crud-credencials/pkg/handler/domain/compliance/benchmark_schema"
	varconfig "projeto-crud-credencials/pkg/handler/domain/core/varconfig"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *varconfig.Handler, bsHandler *benchmarkSchema.Handler) *gin.Engine {
	router := gin.Default()
	SetupVarConfigRoutes(router, h)
	SetupBenchmarkSchemaRoutes(router, bsHandler)

	return router
}
