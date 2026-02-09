package benchmark_schema

import (
	"database/sql"
	svcbenchmark "projeto-crud-credencials/internal/service/domain/compliance/benchmark_schema"
	storagebenchmark "projeto-crud-credencials/internal/storage/dynamodb/benchmark_schema"
	postgresbenchmark "projeto-crud-credencials/internal/storage/postgres/benchmark_schema"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func InitHandler(postgresDB *sql.DB, client *dynamodb.Client, tableName string) *Handler {
	relationalRepo := postgresbenchmark.NewRepository(postgresDB)
	nosqlRepo := storagebenchmark.NewRepository(client, tableName)

	serviceImpl := svcbenchmark.NewService(relationalRepo, nosqlRepo)

	return NewHandler(serviceImpl)
}