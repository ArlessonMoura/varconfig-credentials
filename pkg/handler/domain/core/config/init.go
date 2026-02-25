package config

import (
	// Aliases para evitar confusão entre pacotes de mesmo nome
	service "projeto-crud-credentials/internal/service/domain/core/config"
	benchmarkStorage "projeto-crud-credentials/internal/storage/postgres/benchmark"
	storage "projeto-crud-credentials/internal/storage/postgres/config"

	"gorm.io/gorm"
)

// InitHandler inicializa o handler com todas as dependências (Wiring)
func InitHandler(db *gorm.DB) *Handler {
	// 1. Instancia o repository (Implementação do Driver PostgreSQL com GORM)
	repositoryImpl := storage.NewRepository(db)

	// 1b. Instancia o repository de benchmark para validação cruzada
	benchmarkRepo := benchmarkStorage.NewRepository(db)

	// 2. Instancia o service injetando os repositories
	serviceImpl := service.NewService(repositoryImpl, benchmarkRepo)

	// 3. Retorna o handler com o service injetado
	return NewHandler(serviceImpl)
}
