package varconfig

import (
	"context"
	"fmt"

	models "projeto-crud-credencials/pkg/models/varconfig"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Repository implementa IVarConfigRepository usando DynamoDB SDK v2
type Repository struct {
	client    *dynamodb.Client
	tableName string
}

// NewRepository cria uma nova instância do repository
func NewRepository(client *dynamodb.Client, tableName string) *Repository {
	return &Repository{
		client:    client,
		tableName: tableName,
	}
}

// Create persiste um VarConfigItem (PutItem)
func (r *Repository) Create(ctx context.Context, item *models.VarConfigItem) error {
	// O marshalMap transforma o struct models.VarConfigItem em um map[string]types.AttributeValue
	// Isso respeita as tags `dynamodbav` que definimos no model
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
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

// List busca todos os itens de uma partição (Benchmark)
func (r *Repository) List(ctx context.Context, pk string) ([]*models.VarConfigItem, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query dynamodb: %w", err)
	}

	var items []*models.VarConfigItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal query results: %w", err)
	}

	return items, nil
}


// GetByID busca um item específico pela PK e SK
func (r *Repository) GetByID(ctx context.Context, pk string, sk string) (*models.VarConfigItem, error) {
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

	var item models.VarConfigItem
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return &item, nil
}

// Update atualiza um item existente via chaves compostas
func (r *Repository) Update(ctx context.Context, item *models.VarConfigItem) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item for update: %w", err)
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

// Delete remove um item via chaves compostas
func (r *Repository) Delete(ctx context.Context, pk string, sk string) error {
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
