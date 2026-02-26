package routes

import (
	"github.com/Wizzi-Cloud/restwrapper/handler"
	"github.com/gin-gonic/gin"
)

func SetupVarConfigRoutes(router *gin.Engine, h *handler.Handler) {

	// Agrupamento base para evitar repetição de prefixos
	group := router.Group("/org/:orgId/compliance/variables/:benchmark_id")
	{
		// GET – Listar todos os varConfigs do benchmark
		group.GET("", h.HandleGin)

		// POST – Criar um novo varConfig
		group.POST("", h.HandleGin)

		// Rotas que exigem um ID específico
		specific := group.Group("/:id")
		{
			// GET – Obter um varConfig específico
			specific.GET("", h.HandleGin)

			// PUT ou PATCH – Atualizar um varConfig
			specific.PUT("", h.HandleGin)

			// DELETE – Remover um varConfig
			specific.DELETE("", h.HandleGin)
		}
	}
}
