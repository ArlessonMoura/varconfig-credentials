package varconfig

import (
	// Aliases para evitar confusão entre pacotes de mesmo nome
	service "projeto-crud-credencials/internal/service/domain/core/varconfig"
	storage "projeto-crud-credencials/internal/storage/dynamodb/varconfig"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// InitHandler inicializa o handler com todas as dependências (Wiring)
func InitHandler(client *dynamodb.Client, tableName string) *Handler {
	// 1. Instancia o repository (Implementação do Driver DynamoDB)
	repositoryImpl := storage.NewRepository(client, tableName)

	// 2. Instancia o service injetando o repository
	// O serviceImpl deve implementar a interface IVarConfigService definida neste pacote
	serviceImpl := service.NewService(repositoryImpl)

	// 3. Retorna o handler com o service injetado
	return NewHandler(serviceImpl)
}
