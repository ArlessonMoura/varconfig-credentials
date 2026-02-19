package benchmarkschema

import (
	"database/sql"
	service "projeto-crud-credencials/internal/service/domain/compliance/benchmarkschema"
	dynamodbRepo "projeto-crud-credencials/internal/storage/dynamodb/benchmarkschema"
	postgresRepo "projeto-crud-credencials/internal/storage/postgres/benchmarkschema"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func InitHandler(postgresDB *sql.DB, client *dynamodb.Client, tableName string) *Handler {
	relationalRepo := postgresRepo.NewRepository(postgresDB)
	nosqlRepo := dynamodbRepo.NewRepository(client, tableName)

	serviceImpl := service.NewService(relationalRepo, nosqlRepo)

	return NewHandler(serviceImpl)
}