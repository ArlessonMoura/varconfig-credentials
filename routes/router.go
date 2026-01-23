package routes

import (
	"projeto-crud-credencials/pkg/handler/domain/core/varconfig"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *varconfig.Handler) *gin.Engine {
	router := gin.Default()
	SetupVarConfigRoutes(router, h)

	return router
}
