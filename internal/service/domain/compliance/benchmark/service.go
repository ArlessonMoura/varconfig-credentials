package benchmark

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dto "projeto-crud-credentials/dto/benchmark"
	ports "projeto-crud-credentials/internal/service"
	models "projeto-crud-credentials/pkg/models/benchmark"
)

var (
	ErrInvalidInput   = errors.New("invalid input: name is required")
	ErrEmptySchema    = errors.New("schema cannot be empty")
	ErrSchemaNotFound = errors.New("schema not found")
)

type Service struct {
	relationalRepo ports.IBenchmarkRepository
}

func NewService(relationalRepo ports.IBenchmarkRepository) *Service {
	return &Service{relationalRepo: relationalRepo}
}

func (s *Service) Create(ctx context.Context, name string, schemaRequest *dto.BenchmarkCreateRequestDTO) (dto.BenchmarkCreateResponseDTO, error) {
	// Validações
	if name == "" {
		return dto.BenchmarkCreateResponseDTO{}, ErrInvalidInput
	}

	if len(schemaRequest.Schema) == 0 {
		return dto.BenchmarkCreateResponseDTO{}, ErrEmptySchema
	}

	// Marshal full schema definition into JSON and store in jsonb column
	schemaJSON, err := json.Marshal(schemaRequest.Schema)
	if err != nil {
		return dto.BenchmarkCreateResponseDTO{}, fmt.Errorf("failed to marshal schema: %w", err)
	}

	postgresSchema := &models.BenchmarkSchemaPostgreSQL{
		Name:   name,
		Schema: schemaJSON,
	}

	if err := s.relationalRepo.Create(ctx, postgresSchema); err != nil {
		return dto.BenchmarkCreateResponseDTO{}, fmt.Errorf("failed to create schema in relational db: %w", err)
	}

	return dto.BenchmarkCreateResponseDTO{
		ID:        fmt.Sprintf("%d", postgresSchema.ID),
		Name:      name,
		CreatedAt: postgresSchema.CreatedAt.Format(time.RFC3339),
	}, nil
}

// List recupera todos os schemas disponíveis (versão enxuta sem SchemaBody)
func (s *Service) List(ctx context.Context) (*dto.BenchmarkListResponseDTO, error) {
	items, err := s.relationalRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas from repository: %w", err)
	}

	if items == nil {
		return &dto.BenchmarkListResponseDTO{Data: []dto.BenchmarkResponseDTO{}, Count: 0}, nil
	}

	var data []dto.BenchmarkResponseDTO
	for _, item := range items {
		data = append(data, dto.BenchmarkResponseDTO{
			ID:         fmt.Sprintf("%d", item.ID),
			Name:       item.Name,
			SchemaBody: nil,
			CreatedAt:  item.CreatedAt.Format(time.RFC3339),
		})
	}

	return &dto.BenchmarkListResponseDTO{
		Data:  data,
		Count: len(data),
	}, nil
}

// GetByID busca um schema pelo ID
func (s *Service) GetByID(ctx context.Context, id string) (*dto.BenchmarkResponseDTO, error) {
	// id vem como string (do handler). Converter para int64
	var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	item, err := s.relationalRepo.GetByID(ctx, intID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema from repository: %w", err)
	}
	if item == nil {
		return nil, ErrSchemaNotFound
	}

	// Unmarshal stored JSON into DTO structure and convert to string map for compatibility
	var schemaDef map[string]dto.SchemaFieldDefinition
	if err := json.Unmarshal(item.Schema, &schemaDef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stored schema: %w", err)
	}

	schemaBody := convertSchemaBodyToStringMap(schemaDef)

	return &dto.BenchmarkResponseDTO{
		ID:         fmt.Sprintf("%d", item.ID),
		Name:       item.Name,
		SchemaBody: schemaBody,
		CreatedAt:  item.CreatedAt.Format(time.RFC3339),
	}, nil
}

func convertSchemaBodyToStringMap(schema map[string]dto.SchemaFieldDefinition) map[string]string {
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
