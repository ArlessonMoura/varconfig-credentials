package benchmark_schema

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"projeto-crud-credencials/pkg/models"
)

type Repository struct {
	client    *dynamodb.Client
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) *Repository {
	return &Repository{
		client:    client,
		tableName: tableName,
	}
}

// Gera PK a partir do ID do banco relacional
func (r *Repository) Save(ctx context.Context, item models.BenchmarkSchemaNoSQL) error {
	// Gerar PK e SK a partir do ID
	pk := "SCHEMA#" + item.ID
	sk := "SCHEMA#" + item.ID

	// Criar mapa com os campos incluindo PK e SK
	itemWithKeys := map[string]interface{}{
		"PK":         pk,
		"SK":         sk,
		"ID":         item.ID,
		"schema_body": item.SchemaBody,
		"created_at": item.CreatedAt,
	}

	av, err := attributevalue.MarshalMap(itemWithKeys)
	if err != nil {
		return fmt.Errorf("failed to marshal benchmark schema: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("failed to put item in dynamodb: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (*models.BenchmarkSchemaNoSQL, error) {
	pk := "SCHEMA#" + id
	sk := "SCHEMA#" + id

	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": pk,
		"SK": sk,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item from dynamodb: %w", err)
	}

	if result.Item == nil {
		return nil, nil // Not found não é erro de sistema
	}

	var item models.BenchmarkSchemaNoSQL
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal benchmark schema: %w", err)
	}

	return &item, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	pk := "SCHEMA#" + id
	sk := "SCHEMA#" + id

	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": pk,
		"SK": sk,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal delete key: %w", err)
	}

	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	})

	if err != nil {
		return fmt.Errorf("failed to delete item from dynamodb: %w", err)
	}

	return nil
}
