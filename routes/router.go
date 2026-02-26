package routes

import (
	"github.com/Wizzi-Cloud/restwrapper/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(hConfig *handler.Handler, hBench *handler.Handler) *gin.Engine {
	router := gin.Default()
	SetupVarConfigRoutes(router, hConfig)
	SetupBenchmarkSchemaRoutes(router, hBench)

	return router
}
