package varconfig

import (
	"projeto-crud-credencials/internal/common/logger"
	"projeto-crud-credencials/internal/common/metrics"
	svcvarconfig "projeto-crud-credencials/internal/service/domain/varconfig"
	storagevarconfig "projeto-crud-credencials/internal/storage/dynamodb/varconfig"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// InitHandler inicializa o handler com todas as dependências
func InitHandler(client *dynamodb.Client, tableName string) *Handler {
	// Instancia o repository (storage)
	repositoryImpl := storagevarconfig.NewRepository(client, tableName)

	// Instancia o service com logger e metrics
	loggerImpl := logger.GetGlobalLogger()
	metricsImpl := &metrics.DefaultCollector{}
	serviceImpl := svcvarconfig.NewService(repositoryImpl, loggerImpl, metricsImpl)

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
