package config

import (
	service "projeto-crud-credentials/internal/service/domain/core/config"
	benchmarkStorage "projeto-crud-credentials/internal/storage/postgres/benchmark"
	storage "projeto-crud-credentials/internal/storage/postgres/config"

	"github.com/Wizzi-Cloud/restwrapper/handler"
	"gorm.io/gorm"
)

func InitHandler(db *gorm.DB) *handler.Handler {
	// 1. Instancia o repository (Implementação do Driver PostgreSQL com GORM)
	repositoryImpl := storage.NewRepository(db)

	// 1b. Instancia o repository de benchmark para validação cruzada
	benchmarkRepo := benchmarkStorage.NewRepository(db)

	serviceImpl := service.NewService(repositoryImpl, benchmarkRepo)

	// 2. Cria o handler que implementa IRestHandler
	h := NewHandler(serviceImpl)

	// 3. Retorna o handler universal usando a factory
	return handler.NewHandler(h)
}
