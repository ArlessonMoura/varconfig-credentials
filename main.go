package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	handlervarconfig "projeto-crud-credencials/pkg/handler/domain/core/varconfig"
	"projeto-crud-credencials/routes"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Bootstrap configura as dependências e inicia a aplicação seguindo a arquitetura definida.
func Bootstrap() {
	// 1. Configurar contexto com timeout para a inicialização
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Configurar conexão com DynamoDB via AWS SDK v2
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração AWS: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	
	// O nome da tabela pode vir de uma variável de ambiente conforme boa prática
	tableName := os.Getenv("DYNAMODB_TABLE_VARCONFIG")
	if tableName == "" {
		tableName = "VarConfigs"
	}

	// 3. Criar o Handler (Injeção de Dependência)
	// O InitHandler realiza o wiring interno: Repository -> Service -> Handler
	varConfigHandler := handlervarconfig.InitHandler(client, tableName)

	// 4. Configurar o Router centralizado em /routes
	// O SetupRouter recebe o handler e registra todas as rotas de domínio
	router := routes.SetupRouter(varConfigHandler)

	// 5. Configurar o servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Iniciar servidor seguindo as normas de observabilidade
	fmt.Printf("Servidor iniciado em http://localhost:%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

// main é o ponto de entrada da aplicação
func main() {
	Bootstrap()
}