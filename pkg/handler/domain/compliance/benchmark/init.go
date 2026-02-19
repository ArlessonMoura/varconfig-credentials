package benchmark

import (
	"database/sql"
	service "projeto-crud-credentials/internal/service/domain/compliance/benchmark"
	dynamodbRepo "projeto-crud-credentials/internal/storage/dynamodb/benchmark"
	postgresRepo "projeto-crud-credentials/internal/storage/postgres/benchmark"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func InitHandler(postgresDB *sql.DB, client *dynamodb.Client, tableName string) *Handler {
	relationalRepo := postgresRepo.NewRepository(postgresDB)
	nosqlRepo := dynamodbRepo.NewRepository(client, tableName)

	serviceImpl := service.NewService(relationalRepo, nosqlRepo)

	return NewHandler(serviceImpl)
}