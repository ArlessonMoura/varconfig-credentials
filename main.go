package main

import (
	"context"
	"fmt"

	handlervarconfig "projeto-crud-credencials/pkg/handler/varconfig"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// ExampleMain demonstra como integrar o CRUD de VarConfig em sua aplicação
func Bootstrap() {
	// 1. Configurar conexão com DynamoDB
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Erro ao carregar configuração AWS:", err)
		return
	}

	client := dynamodb.NewFromConfig(cfg)
	tableName := "VarConfigs" // Nome da tabela no DynamoDB

	// 2. Criar handler
	handler := handlervarconfig.InitHandler(client, tableName)

	// 3. Configurar router Gin
	router := gin.Default()

	// 4. Registrar rotas de VarConfig
	handlervarconfig.RegisterRoutes(router, handler)

	// 5. Iniciar servidor
	router.Run(":8080")
}

// main é o ponto de entrada da aplicação
func main() {
	Bootstrap()
}