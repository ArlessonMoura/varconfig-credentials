package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	handlerbenchmarkschema "projeto-crud-credentials/pkg/handler/domain/compliance/benchmarkschema"
	handlervarconfig "projeto-crud-credentials/pkg/handler/domain/core/varconfig"
	"projeto-crud-credentials/routes"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	_ "github.com/lib/pq"
)

func Bootstrap() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Conectar ao PostgreSQL
	postgresConnection := os.Getenv("POSTGRES_CONNECTION_STRING")
	if postgresConnection == "" {
		log.Fatalf("Erro: Variável de ambiente POSTGRES_CONNECTION_STRING não configurada")
	}

	sqlDB, err := sql.Open("postgres", postgresConnection)
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("Erro ao validar conexão PostgreSQL: %v", err)
	}

	// 2. Conectar ao DynamoDB
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração AWS: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	tableName := os.Getenv("DYNAMODB_TABLE_VARCONFIG")
	if tableName == "" {
		log.Fatal("Erro: Variável de ambiente DYNAMODB_TABLE_VARCONFIG nao configurada")
	}

	benchmarkSchemaTableName := os.Getenv("DYNAMODB_TABLE_BENCHMARK_SCHEMA")
	if benchmarkSchemaTableName == "" {
		log.Fatal("Erro: Variável de ambiente DYNAMODB_TABLE_BENCHMARK_SCHEMA nao configurada")
	}

	varConfigHandler := handlervarconfig.InitHandler(client, tableName)
	benchmarkSchemaHandler := handlerbenchmarkschema.InitHandler(sqlDB, client, benchmarkSchemaTableName)

	router := routes.SetupRouter(varConfigHandler, benchmarkSchemaHandler)

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

	fmt.Printf("Servidor iniciado em http://localhost:%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func main() {
	Bootstrap()
}
