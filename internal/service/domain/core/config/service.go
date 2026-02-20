package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dtoBench "projeto-crud-credentials/dto/benchmark"
	dto "projeto-crud-credentials/dto/config"
	ports "projeto-crud-credentials/internal/service"
	models "projeto-crud-credentials/pkg/models/config"
)

// Erros de domínio básicos
var (
	ErrNotFound          = errors.New("varconfig not found")
	ErrInvalidInput      = errors.New("invalid input: orgID and benchmarkID are required")
	ErrPayloadValidation = errors.New("payload validation failed")
)

type Service struct {
	repository    ports.IVarConfigRepository
	benchmarkRepo ports.IBenchmarkRepository
}

func NewService(repository ports.IVarConfigRepository, benchmarkRepo ports.IBenchmarkRepository) *Service {
	return &Service{
		repository:    repository,
		benchmarkRepo: benchmarkRepo,
	}
}

// Create cria um novo VarConfig com ID gerado automaticamente pelo PostgreSQL
func (s *Service) Create(ctx context.Context, orgID string, benchmarkID string, input *dto.CreateVarConfigRequest) (*dto.VarConfigResponse, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, ErrInvalidInput
	}

	// Validar payload contra o schema do benchmark (checagem mínima de tipos)
	if err := s.validatePayloadAgainstBenchmark(ctx, benchmarkID, input.Payload); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	item, err := s.repository.Create(ctx, orgID, benchmarkID, input.Payload)
	if err != nil {
		return nil, fmt.Errorf("storage error: %w", err)
	}

	return s.mapItemToResponse(item), nil
}

// List retorna todas as configurações de benchmark (versão leve sem payload)
func (s *Service) List(ctx context.Context, orgID, benchmarkID string) (*dto.ListVarConfigsResponse, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, ErrInvalidInput
	}

	items, err := s.repository.List(ctx, orgID, benchmarkID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	var data []dto.ListVarConfigResponse
	for _, item := range items {
		data = append(data, dto.ListVarConfigResponse{
			ID:          fmt.Sprintf("%d", item.ID),
			OrgID:       item.OrgID,
			BenchmarkID: item.BenchmarkID,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &dto.ListVarConfigsResponse{Data: data}, nil
}

// GetByID retorna todas as configurações de um benchmark específico
func (s *Service) GetByID(ctx context.Context, orgID, benchmarkID, id string) (*dto.VarConfigResponse, error) {
	if orgID == "" || benchmarkID == "" || id == "" {
		return nil, ErrInvalidInput
	}

	// Converter ID string para int64
	var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	item, err := s.repository.GetByID(ctx, intID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if item == nil {
		return nil, ErrNotFound
	}

	// Validar que o item pertence aos IDs fornecidos
	if item.OrgID != orgID || item.BenchmarkID != benchmarkID {
		return nil, ErrNotFound
	}

	return s.mapItemToResponse(item), nil
}

func (s *Service) Update(ctx context.Context, orgID, benchmarkID, id string, input *dto.UpdateVarConfigRequest) (*dto.VarConfigResponse, error) {
	if orgID == "" || benchmarkID == "" || id == "" {
		return nil, ErrInvalidInput
	}

	// Converter ID string para int64
	var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// Verificar que o item existe e pertence aos IDs fornecidos
	existing, err := s.repository.GetByID(ctx, intID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	if existing.OrgID != orgID || existing.BenchmarkID != benchmarkID {
		return nil, ErrNotFound
	}

	// Validar payload contra o schema do benchmark (minima de checagem de tipos)
	if err := s.validatePayloadAgainstBenchmark(ctx, benchmarkID, input.Payload); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Atualizar payload
	item, err := s.repository.Update(ctx, intID, input.Payload)
	if err != nil {
		return nil, fmt.Errorf("update error: %w", err)
	}

	return s.mapItemToResponse(item), nil
}

func (s *Service) Delete(ctx context.Context, orgID, benchmarkID, id string) error {
	if orgID == "" || benchmarkID == "" || id == "" {
		return ErrInvalidInput
	}

	// Converter ID string para int64
	var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	// Verificar que o item existe e pertence aos IDs fornecidos
	existing, err := s.repository.GetByID(ctx, intID)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}
	if existing == nil {
		return ErrNotFound
	}

	if existing.OrgID != orgID || existing.BenchmarkID != benchmarkID {
		return ErrNotFound
	}

	return s.repository.Delete(ctx, intID)
}

// mapItemToResponse converte VarConfigPostgreSQL para VarConfigResponse
func (s *Service) mapItemToResponse(item *models.VarConfigPostgreSQL) *dto.VarConfigResponse {
	var payload map[string]any
	if item.Payload != nil {
		if err := json.Unmarshal(item.Payload, &payload); err != nil {
			payload = nil
		}
	}

	return &dto.VarConfigResponse{
		ID:          fmt.Sprintf("%d", item.ID),
		OrgID:       item.OrgID,
		BenchmarkID: item.BenchmarkID,
		Payload:     payload,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	}
}

// validatePayloadAgainstBenchmark faz uma checagem mínima de tipos do payload
func (s *Service) validatePayloadAgainstBenchmark(ctx context.Context, benchmarkID string, payload map[string]any) error {
	// converter benchmarkID para int64 (o repositório trabalha com int64)
	var intID int64
	if _, err := fmt.Sscan(benchmarkID, &intID); err != nil {
		return fmt.Errorf("invalid benchmark id format: %w", err)
	}

	// Buscar schema
	bench, err := s.benchmarkRepo.GetByID(ctx, intID)
	if err != nil {
		return fmt.Errorf("failed to fetch benchmark schema: %w", err)
	}
	if bench == nil {
		return fmt.Errorf("benchmark schema not found")
	}

	var schemaDef map[string]dtoBench.SchemaFieldDefinition
	if err := json.Unmarshal(bench.Schema, &schemaDef); err != nil {
		return fmt.Errorf("failed to unmarshal benchmark schema: %w", err)
	}

	// Validar cada campo requerido e tipo básico
	for fieldName, def := range schemaDef {
		val, exists := payload[fieldName]
		if def.Required && !exists {
			return fmt.Errorf("field '%s' is required: %w", fieldName, ErrPayloadValidation)
		}
		if exists {
			switch def.Type {
			case dtoBench.FieldTypeString:
				if _, ok := val.(string); !ok {
					return fmt.Errorf("field '%s' must be string: %w", fieldName, ErrPayloadValidation)
				}
			case dtoBench.FieldTypeNumber:
				switch val.(type) {
				case float64, float32, int, int64, int32:
					// ok
				default:
					return fmt.Errorf("field '%s' must be number: %w", fieldName, ErrPayloadValidation)
				}
			case dtoBench.FieldTypeBoolean:
				if _, ok := val.(bool); !ok {
					return fmt.Errorf("field '%s' must be boolean: %w", fieldName, ErrPayloadValidation)
				}
			case dtoBench.FieldTypeArray:
				// ensure slice
				arr, ok := val.([]any)
				if !ok {
					return fmt.Errorf("field '%s' must be array: %w", fieldName, ErrPayloadValidation)
				}
				if def.Items != nil {
					for _, elem := range arr {
						switch def.Items.Type {
						case dtoBench.FieldTypeString:
							if _, ok := elem.(string); !ok {
								return fmt.Errorf("array field '%s' elements must be string: %w", fieldName, ErrPayloadValidation)
							}
						case dtoBench.FieldTypeNumber:
							switch elem.(type) {
							case float64, float32, int, int64, int32:
							default:
								return fmt.Errorf("array field '%s' elements must be number: %w", fieldName, ErrPayloadValidation)
							}
						case dtoBench.FieldTypeBoolean:
							if _, ok := elem.(bool); !ok {
								return fmt.Errorf("array field '%s' elements must be boolean: %w", fieldName, ErrPayloadValidation)
							}
						}
					}
				}
			}
		}
	}

	return nil
}
