// Package varconfig provides DynamoDB repository implementation for VarConfig storage.
//
// This package implements the VarConfigRepository interface using AWS DynamoDB
// with optimized access patterns using composite primary keys (PK+SK).
//
// Key structure:
//
//	PK = "ORG#{orgId}#BENCH#{benchmark_id}"
//	SK = "VARCONFIG#{id}"
//
// All operations use efficient GetItem/Query patterns, avoiding Scan operations.
package varconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	svcvarconfig "projeto-crud-credencials/internal/service/domain/varconfig"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Repository implementa VarConfigRepository usando DynamoDB
type Repository struct {
	client *dynamodb.Client
	tableName string
}

// NewRepository cria uma nova instância do repository
func NewRepository(client *dynamodb.Client, tableName string) *Repository {
	return &Repository{
		client: client,
		tableName: tableName,
	}
}

// DynamoDBModel é o modelo usado internamente pelo storage
type DynamoDBModel struct {
	PK          string `dynamodbav:"PK"`           // "ORG#{orgId}#BENCH#{benchmark_id}"
	SK          string `dynamodbav:"SK"`           // "VARCONFIG#{id}"
	ID          string `dynamodbav:"id"`
	OrgID       string `dynamodbav:"org_id"`
	BenchmarkID string `dynamodbav:"benchmark_id"`
	Payload     string `dynamodbav:"payload"`
	CreatedAt   string `dynamodbav:"created_at"`
	UpdatedAt   string `dynamodbav:"updated_at"`
}

// Helper functions para criar PK e SK
func buildPK(orgID int64, benchmarkID string) string {
	return fmt.Sprintf("ORG#%d#BENCH#%s", orgID, benchmarkID)
}

func buildSK(id int64) string {
	return fmt.Sprintf("VARCONFIG#%d", id)
}

// Save persiste um novo VarConfig
func (r *Repository) Save(config svcvarconfig.VarConfig) (svcvarconfig.VarConfig, error) {
	now := time.Now().UTC()
	
	// Gerar ID único se não existir
	if config.ID == 0 {
		config.ID = time.Now().UnixNano()
	}

	model := DynamoDBModel{
		PK:          buildPK(config.OrgID, config.BenchmarkID),
		SK:          buildSK(config.ID),
		ID:          fmt.Sprintf("%d", config.ID),
		OrgID:       fmt.Sprintf("%d", config.OrgID),
		BenchmarkID: config.BenchmarkID,
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}

	// Serializar payload
	if config.Payload != nil {
		payloadData, err := json.Marshal(config.Payload)
		if err != nil {
			return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao serializar payload: %w", err)
		}
		model.Payload = string(payloadData)
	}

	// Converter para formato DynamoDB
	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao converter item para DynamoDB: %w", err)
	}

	// Inserir no DynamoDB
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao salvar VarConfig: %w", err)
	}

	// Definir timestamps no objeto de retorno
	config.CreatedAt = now
	config.UpdatedAt = now

	return config, nil
}

// FindByID obtém um VarConfig específico
func (r *Repository) FindByID(orgID int64, benchmarkID string, id int64) (svcvarconfig.VarConfig, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": buildPK(orgID, benchmarkID),
		"SK": buildSK(id),
	})
	if err != nil {
		return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao criar chave: %w", err)
	}

	result, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	})

	if err != nil {
		return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao buscar VarConfig: %w", err)
	}

	if result.Item == nil {
		return svcvarconfig.VarConfig{}, fmt.Errorf("VarConfig não encontrado")
	}

	var model DynamoDBModel
	err = attributevalue.UnmarshalMap(result.Item, &model)
	if err != nil {
		return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao desserializar item: %w", err)
	}

	return r.fromModel(&model)
}

// FindAllByBenchmark lista todos os VarConfigs de um benchmark
func (r *Repository) FindAllByBenchmark(orgID int64, benchmarkID string) ([]svcvarconfig.VarConfig, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: buildPK(orgID, benchmarkID)},
		},
	}

	result, err := r.client.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar VarConfigs: %w", err)
	}

	var configs []svcvarconfig.VarConfig
	for _, item := range result.Items {
		var model DynamoDBModel
		err := attributevalue.UnmarshalMap(item, &model)
		if err != nil {
			continue // Pular itens inválidos
		}

		config, err := r.fromModel(&model)
		if err != nil {
			continue // Pular itens com erro
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// Update atualiza um VarConfig existente
func (r *Repository) Update(config svcvarconfig.VarConfig) (svcvarconfig.VarConfig, error) {
	now := time.Now().UTC()
	config.UpdatedAt = now

	model := DynamoDBModel{
		PK:          buildPK(config.OrgID, config.BenchmarkID),
		SK:          buildSK(config.ID),
		ID:          fmt.Sprintf("%d", config.ID),
		OrgID:       fmt.Sprintf("%d", config.OrgID),
		BenchmarkID: config.BenchmarkID,
		UpdatedAt:   now.Format(time.RFC3339),
	}

	// Serializar payload
	if config.Payload != nil {
		payloadData, err := json.Marshal(config.Payload)
		if err != nil {
			return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao serializar payload: %w", err)
		}
		model.Payload = string(payloadData)
	}

	// Construir expression de atualização
	updateExpression := "SET payload = :payload, updated_at = :updated_at"
	expressionValues := map[string]types.AttributeValue{
		":payload":     &types.AttributeValueMemberS{Value: model.Payload},
		":updated_at":  &types.AttributeValueMemberS{Value: model.UpdatedAt},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: model.PK},
			"SK": &types.AttributeValueMemberS{Value: model.SK},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
	}

	_, err := r.client.UpdateItem(context.TODO(), input)
	if err != nil {
		return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao atualizar VarConfig: %w", err)
	}

	return config, nil
}

// Delete remove um VarConfig
func (r *Repository) Delete(orgID int64, benchmarkID string, id int64) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: buildPK(orgID, benchmarkID)},
			"SK": &types.AttributeValueMemberS{Value: buildSK(id)},
		},
	}

	_, err := r.client.DeleteItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("erro ao deletar VarConfig: %w", err)
	}

	return nil
}

// fromModel converte DynamoDBModel para VarConfig
func (r *Repository) fromModel(model *DynamoDBModel) (svcvarconfig.VarConfig, error) {
	config := svcvarconfig.VarConfig{
		Payload: make(map[string]any),
	}

	// Parse ID
	id, err := parseID(model.ID)
	if err == nil {
		config.ID = id
	}

	// Parse OrgID
	orgID, err := parseID(model.OrgID)
	if err == nil {
		config.OrgID = orgID
	}

	config.BenchmarkID = model.BenchmarkID

	// Parse timestamps
	if createdAt, err := time.Parse(time.RFC3339, model.CreatedAt); err == nil {
		config.CreatedAt = createdAt
	}

	if updatedAt, err := time.Parse(time.RFC3339, model.UpdatedAt); err == nil {
		config.UpdatedAt = updatedAt
	}

	// Parse payload
	if model.Payload != "" {
		err := json.Unmarshal([]byte(model.Payload), &config.Payload)
		if err != nil {
			return svcvarconfig.VarConfig{}, fmt.Errorf("erro ao desserializar payload: %w", err)
		}
	}

	return config, nil
}

// parseID converte string ID para int64
func parseID(idStr string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	return id, err
}
