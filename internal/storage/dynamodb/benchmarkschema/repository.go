package benchmarkschema

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	models "projeto-crud-credentials/pkg/models/benchmarkschema"
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

// Create armazena o schema no DynamoDB usando o ID como chave de partição
func (r *Repository) Create(ctx context.Context, item *models.BenchmarkSchemaDynamoDB) error {
	// Gerar PK a partir do ID (sem SK redundante)
	pk := "SCHEMA#" + item.ID

	// Criar mapa com os campos incluindo PK
	itemWithKeys := map[string]interface{}{
		"PK":          pk,
		"ID":          item.ID,
		"schema_body": item.SchemaBody,
		"created_at":  item.CreatedAt,
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

// List retorna apenas campos essenciais para listagem (ID, CreatedAt)
func (r *Repository) List(ctx context.Context) ([]*models.BenchmarkSchemaListDynamoDB, error) {
	// Query pela partição SCHEMA# (todos os schemas começam com SCHEMA#)
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("begins_with(PK, :pk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk_prefix": &types.AttributeValueMemberS{Value: "SCHEMA#"},
		},
		ProjectionExpression: aws.String("ID, created_at"),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query benchmark schemas from dynamodb: %w", err)
	}

	var items []models.BenchmarkSchemaListDynamoDB
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal benchmark schemas: %w", err)
	}

	// Converter slice de structs para slice de ponteiros
	var pointerItems []*models.BenchmarkSchemaListDynamoDB
	for i := range items {
		pointerItems = append(pointerItems, &items[i])
	}

	return pointerItems, nil
}

// GetByID busca um schema pelo ID
func (r *Repository) GetByID(ctx context.Context, id string) (*models.BenchmarkSchemaDynamoDB, error) {
	pk := "SCHEMA#" + id

	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": pk,
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

	var item models.BenchmarkSchemaDynamoDB
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal benchmark schema: %w", err)
	}

	return &item, nil
}

//====================
//
//====================
	
// Update atualiza um schema existente no DynamoDB
func (r *Repository) Update(ctx context.Context, item *models.BenchmarkSchemaDynamoDB) error {
	pk := "SCHEMA#" + item.ID

	// Criar mapa com os campos incluindo PK
	itemWithKeys := map[string]interface{}{
		"PK":          pk,
		"ID":          item.ID,
		"schema_body": item.SchemaBody,
		"created_at":  item.CreatedAt,
	}

	av, err := attributevalue.MarshalMap(itemWithKeys)
	if err != nil {
		return fmt.Errorf("failed to marshal benchmark schema for update: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("failed to update item in dynamodb: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	pk := "SCHEMA#" + id

	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": pk,
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
