package varconfig

import (
	svcvarconfig "projeto-crud-credencials/internal/service/domain/varconfig"
	storagevarconfig "projeto-crud-credencials/internal/storage/dynamodb/varconfig"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// InitHandler inicializa o handler com todas as dependências
// Esta função implementa o padrão de wiring descrito em ARCHITECTURE.md
func InitHandler(client *dynamodb.Client, tableName string) *Handler {
	// Instancia o repository (storage)
	repositoryImpl := storagevarconfig.NewRepository(client, tableName)

	// Instancia o service
	serviceImpl := svcvarconfig.NewService(repositoryImpl)

	// Retorna o handler com o service injetado
	return NewHandler(serviceImpl)
}

// RegisterRoutes registra as rotas de VarConfig no router
func RegisterRoutes(router *gin.Engine, handler *Handler) {
	varconfigRoutes := router.Group("/org/:orgId/compliance/variables")
	{
		benchmarkGroup := varconfigRoutes.Group("/:benchmark_id")
		{
			benchmarkGroup.GET("", handler.ListByBenchmark)
			benchmarkGroup.POST("", handler.Create)

			idGroup := benchmarkGroup.Group("/:id")
			{
				idGroup.GET("", handler.GetByID)
				idGroup.PUT("", handler.Update)
				idGroup.PATCH("", handler.Update)
				idGroup.DELETE("", handler.Delete)
			}
		}
	}
}
