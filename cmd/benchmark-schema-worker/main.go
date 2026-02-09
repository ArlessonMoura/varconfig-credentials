package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	benchmarks "projeto-crud-credencials/internal/service/domain/compliance/benchmark_schema"
	dynamodbrepo "projeto-crud-credencials/internal/storage/dynamodb/benchmark_schema"
	postgresrepo "projeto-crud-credencials/internal/storage/postgres/benchmark_schema"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	_ "github.com/lib/pq"
)

func Bootstrap() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Conectar ao PostgreSQL (para registrar IDs)
	postgresConnection := os.Getenv("POSTGRES_CONNECTION_STRING")
	if postgresConnection == "" {
		log.Fatalf("Erro: Variável de ambiente POSTGRES_CONNECTION_STRING não configurada")
	}

	sqlDB, err := sql.Open("postgres", postgresConnection)
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("Erro ao validar conexão PostgreSQL: %v", err)
	}

	postgresRepo := postgresrepo.NewRepository(sqlDB)

	// 2. Conectar ao DynamoDB (para armazenar JSON Schema)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração AWS: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("DYNAMODB_TABLE_BENCHMARK_SCHEMA")
	if tableName == "" {
		log.Fatalf("Erro: Variável de ambiente DYNAMODB_TABLE_BENCHMARK_SCHEMA não configurada")
	}

	dynamodbRepo := dynamodbrepo.NewRepository(client, tableName)


	// Não entendi onde sera consumido esse serviço, mas ele precisa existir para injetado em algum package de handler, então vou criar ele aqui mesmo. O importante é que ele seja inicializado com os repositórios corretos.	
	_ = benchmarks.NewService(postgresRepo, dynamodbRepo) 
	// Service será instanciado pelos processos que usarem esse pacote
}

func main() {
	Bootstrap()
}