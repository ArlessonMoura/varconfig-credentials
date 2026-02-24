package benchmark

import (
	service "projeto-crud-credentials/internal/service/domain/compliance/benchmark"
	postgresRepo "projeto-crud-credentials/internal/storage/postgres/benchmark"

	"gorm.io/gorm"
)

func InitHandler(gormDB *gorm.DB) *Handler {
	relationalRepo := postgresRepo.NewRepository(gormDB)

	serviceImpl := service.NewService(relationalRepo)

	return NewHandler(serviceImpl)
}
