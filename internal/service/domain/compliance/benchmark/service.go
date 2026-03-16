package benchmark

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dto "projeto-crud-credentials/dto/benchmark"
	ports "projeto-crud-credentials/internal/service"
	"projeto-crud-credentials/pkg/models"
)

var (
	ErrInvalidInput     = errors.New("invalid input: name is required")
	ErrEmptySchema      = errors.New("schema cannot be empty")
	ErrSchemaNotFound   = errors.New("schema not found")
	ErrDuplicateName    = errors.New("schema with this name already exists")
)

type Service struct {
	relationalRepo ports.IBenchmarkRepository
}

func NewService(relationalRepo ports.IBenchmarkRepository) *Service {
	return &Service{relationalRepo: relationalRepo}
}

func (s *Service) Create(ctx context.Context, schemaRequest *dto.BenchmarkCreateRequestDTO) (*models.BenchmarkSchema, error) {
	
	if err := schemaRequest.Validate(); err != nil {
		return nil, err
	}

	if len(schemaRequest.Schema) == 0 {
		return nil, ErrEmptySchema
	}

	schemaJSON, err := json.Marshal(schemaRequest.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	schema := &models.BenchmarkSchema{
		Name:    schemaRequest.Name,
		Version: schemaRequest.Version,
		Schema:  schemaJSON,
	}

	if err := s.relationalRepo.Create(ctx, schema); err != nil {
		// Tratamento específico para violação de unique constraint
		if isUniqueViolationError(err) {
			return nil, ErrDuplicateName
		}
		return nil, fmt.Errorf("failed to create schema in relational db: %w", err)
	}

	return schema, nil
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
			Version:    item.Version,
			SchemaBody: nil,
			CreatedAt:  item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  item.UpdatedAt.Format(time.RFC3339),
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
		Version:    item.Version,
		SchemaBody: schemaBody,
		CreatedAt:  item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  item.UpdatedAt.Format(time.RFC3339),
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

// isUniqueViolationError verifica se o erro é uma violação de constraint unique
func isUniqueViolationError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "unique constraint") || contains(errStr, "duplicate key") || contains(errStr, "UNIQUE violation")
}

// contains verifica se uma substring existe em uma string (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
				s[len(s)-len(substr):] == substr || 
				indexOf(s, substr) >= 0)))
}

// indexOf retorna o índice da primeira ocorrência de substr em s
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Delete remove um benchmark schema e automaticamente todas as configurações associadas (CASCADE)
func (s *Service) Delete(ctx context.Context, id string) error {
	// Converter ID string para int64
	var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	// Verificar se o schema existe antes de deletar
	_, err := s.relationalRepo.GetByID(ctx, intID)
	if err != nil {
		return fmt.Errorf("failed to get benchmark schema: %w", err)
	}

	// Deletar o schema (cascade deletará automaticamente os VarConfigs associados)
	if err := s.relationalRepo.Delete(ctx, &intID); err != nil {
		return fmt.Errorf("failed to delete benchmark schema: %w", err)
	}

	return nil
}
