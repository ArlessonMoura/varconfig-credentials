// Package benchmark_schema provides services for managing benchmark schemas in both relational and NoSQL databases.
package benchmark_schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	dto "projeto-crud-credencials/dto/benchmark_schema"
	repository "projeto-crud-credencials/internal/service"
	models "projeto-crud-credencials/pkg/models/benchmark_schema"
)

var (
	ErrInvalidInput   = errors.New("invalid input: name is required")
	ErrEmptySchema    = errors.New("schema cannot be empty")
	ErrSchemaNotFound = errors.New("schema not found")
)

type Service struct {
	relationalRepo repository.IRelationalRepository
	nosqlRepo      repository.INoSQLRepository

}

func NewService(relRepo repository.IRelationalRepository, noSQLRepo repository.INoSQLRepository) *Service {
	return &Service{
		relationalRepo: relRepo,
		nosqlRepo:      noSQLRepo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	name string,
	schemaRequest *dto.InternalRegisterSchemaRequest,
) (dto.RegisterSchemaResponse, error) {
	// Validações
	if name == "" {
		return dto.RegisterSchemaResponse{}, ErrInvalidInput
	}

	if len(schemaRequest.Schema) == 0 {
		return dto.RegisterSchemaResponse{}, ErrEmptySchema
	}

	// 1. Registrar no banco relacional (PostgreSQL)
	relationalModel := &models.BenchmarkSchemaRelational{
		Name: name,
	}

	if err := s.relationalRepo.Create(ctx, relationalModel); err != nil {
		return dto.RegisterSchemaResponse{}, fmt.Errorf("failed to create benchmark schema in relational database: %w", err)
	}

	noSQLModel := &models.BenchmarkSchemaNoSQL{
		ID:         strconv.FormatInt(relationalModel.ID, 10),
		SchemaBody: convertSchemaBodyToStringMap(schemaRequest.Schema),
		CreatedAt:  relationalModel.CreatedAt.Format(time.RFC3339),
	}

	if err := s.nosqlRepo.Create(ctx, noSQLModel); err != nil {
		// 3. Rollback Compensatório: Deleta do SQL se o DynamoDB falhar
		deleteErr := s.relationalRepo.Delete(ctx, relationalModel.ID)
		if deleteErr != nil {
			fmt.Printf("CRITICAL: Failed to rollback schema %d after DynamoDB failure: %v\n", relationalModel.ID, deleteErr)
			return dto.RegisterSchemaResponse{}, fmt.Errorf("failed to save schema in DynamoDB and failed to rollback: %w, rollback error: %w", err, deleteErr)
		}
		return dto.RegisterSchemaResponse{}, fmt.Errorf("failed to save schema in NoSQL database, rolled back relational entry: %w", err)
	}

	return dto.RegisterSchemaResponse{
		ID:        noSQLModel.ID,
		Name:      relationalModel.Name,
		CreatedAt: relationalModel.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetSchemaByID recupera um schema específico pelo ID
func (s *Service) GeByID(ctx context.Context, id string) (dto.BenchmarkSchemaResponse, error) {
	if id == "" {
		return dto.BenchmarkSchemaResponse{}, ErrSchemaNotFound
	}

	item, err := s.nosqlRepo.GetByID(ctx, id)
	if err != nil {
		return dto.BenchmarkSchemaResponse{}, fmt.Errorf("failed to get schema from repository: %w", err)
	}

	if item == nil {
		return dto.BenchmarkSchemaResponse{}, ErrSchemaNotFound
	}

	return dto.BenchmarkSchemaResponse{
		ID:         item.ID,
		Name:       "", // Nome não está disponível no NoSQL, poderia buscar do relacional se necessário
		SchemaBody: item.SchemaBody,
		CreatedAt:  item.CreatedAt,
	}, nil
}

// ListAllSchemas recupera todos os schemas disponíveis
func (s *Service) List(ctx context.Context) (dto.ListBenchmarkSchemasResponse, error) {
	items, err := s.nosqlRepo.List(ctx)
	if err != nil {
		return dto.ListBenchmarkSchemasResponse{}, fmt.Errorf("failed to list schemas from repository: %w", err)
	}

	if items == nil {
		return dto.ListBenchmarkSchemasResponse{Data: []dto.BenchmarkSchemaResponse{}, Count: 0}, nil
	}

	var data []dto.BenchmarkSchemaResponse
	for _, item := range items {
		data = append(data, dto.BenchmarkSchemaResponse{
			ID:         item.ID,
			Name:       "", // Nome não está disponível no NoSQL
			SchemaBody: item.SchemaBody,
			CreatedAt:  item.CreatedAt,
		})
	}

	return dto.ListBenchmarkSchemasResponse{
		Data:  data,
		Count: len(data),
	}, nil
}

func convertSchemaBodyToStringMap(schema map[string]dto.SchemaProperty) map[string]string {
	result := make(map[string]string)
	for k, v := range schema {
		// transforma a struct SchemaProperty em uma string JSON válida
		jsonData, err := json.Marshal(v)
		if err != nil {
			result[k] = fmt.Sprintf(`{"type":"%s","error":"marshal_failed"}`, v.Type)
			continue
		}
		result[k] = string(jsonData)
	}
	return result
}
