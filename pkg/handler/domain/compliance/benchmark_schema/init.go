package benchmark_schema

import (
	"database/sql"
	service "projeto-crud-credencials/internal/service/domain/compliance/benchmark_schema"
	dynamodbRepo "projeto-crud-credencials/internal/storage/dynamodb/benchmark_schema"
	postgresRepo "projeto-crud-credencials/internal/storage/postgres/benchmark_schema"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func InitHandler(postgresDB *sql.DB, client *dynamodb.Client, tableName string) *Handler {
	relationalRepo := postgresRepo.NewRepository(postgresDB)
	nosqlRepo := dynamodbRepo.NewRepository(client, tableName)

	serviceImpl := service.NewService(relationalRepo, nosqlRepo)

	return NewHandler(serviceImpl)
}