package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	hbenchmark "projeto-crud-credentials/pkg/handler/domain/compliance/benchmark"
	hconfig "projeto-crud-credentials/pkg/handler/domain/core/config"
	"projeto-crud-credentials/routes"

	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"projeto-crud-credentials/pkg/models"
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

	// Inicializa GORM usando a conexão sql.DB
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao inicializar GORM: %v", err)
	}

	// AutoMigrate das tabelas essenciais (benchmark schemas e var_configs)
	if err := gormDB.AutoMigrate(&models.BenchmarkSchemaPostgreSQL{}, &models.VarConfigPostgreSQL{}); err != nil {
		log.Fatalf("Erro ao executar AutoMigrate: %v", err)
	}

	varConfigHandler := hconfig.InitHandler(gormDB)
	benchmarkSchemaHandler := hbenchmark.InitHandler(gormDB)

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
