package routes

import (
	"projeto-crud-credencials/pkg/handler/domain/core/varconfig"

	"github.com/gin-gonic/gin"
)

// SetupVarConfigRoutes registra as rotas de VarConfig seguindo a especificação do projeto
func SetupVarConfigRoutes(router *gin.Engine, h *varconfig.Handler) {

	// Agrupamento base para evitar repetição de prefixos
	group := router.Group("/org/:orgId/compliance/variables/:benchmark_id")
	{
		// GET – Listar todos os varConfigs do benchmark
		group.GET("", h.List)

		// POST – Criar um novo varConfig
		group.POST("", h.Create)

		// Rotas que exigem um ID específico
		specific := group.Group("/:id")
		{
			// GET – Obter um varConfig específico
			specific.GET("", h.GetByID)

			// PUT ou PATCH – Atualizar um varConfig
			specific.PUT("", h.Update)

			// DELETE – Remover um varConfig
			specific.DELETE("", h.Delete)
		}
	}
}