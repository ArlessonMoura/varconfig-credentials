package routes

import (
	"projeto-crud-credencials/pkg/handler/domain/compliance/benchmarkschema"
	"projeto-crud-credencials/pkg/handler/domain/core/varconfig"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *varconfig.Handler, bsHandler *benchmarkschema.Handler) *gin.Engine {
	router := gin.Default()
	SetupVarConfigRoutes(router, h)
	SetupBenchmarkSchemaRoutes(router, bsHandler)

	return router
}
