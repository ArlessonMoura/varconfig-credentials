package benchmark

import (
	service "projeto-crud-credentials/internal/service/domain/compliance/benchmark"
	postgresRepo "projeto-crud-credentials/internal/storage/postgres/benchmark"

	"github.com/Wizzi-Cloud/restwrapper/handler"
	"gorm.io/gorm"
)

func InitHandler(gormDB *gorm.DB) *handler.Handler {
	relationalRepo := postgresRepo.NewRepository(gormDB)

	serviceImpl := service.NewService(relationalRepo)

	// Criar o handler que implementa IRestHandler
	h := NewHandler(serviceImpl)

	// Retornar o handler universal usando a factory
	return handler.NewHandler(h)
}
