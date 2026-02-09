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


func Bootstrap() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração AWS: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	
	tableName := os.Getenv("DYNAMODB_TABLE_VARCONFIG")
	if tableName == "" {
		log.Fatal("Erro: Variável de ambiente DYNAMODB_TABLE_VARCONFIG nao configurada")
	}
	
	varConfigHandler := handlervarconfig.InitHandler(client, tableName)
	
	router := routes.SetupRouter(varConfigHandler)

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
