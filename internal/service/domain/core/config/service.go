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
	"projeto-crud-credentials/pkg/models"
)

// Erros de domínio básicos
var (
	ErrNotFound          = errors.New("varconfig not found")
	ErrInvalidInput      = errors.New("invalid input: orgID and benchmarkID are required")
	ErrPayloadValidation = errors.New("payload validation failed")
	ErrDuplicateName     = errors.New("configuration with this name already exists for this organization and benchmark")
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
func (s *Service) Create(ctx context.Context, orgID string, benchmarkID string, input *dto.ConfigCreateRequestDTO) (*dto.ConfigResponseDTO, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, ErrInvalidInput
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Validar payload contra o schema do benchmark (checagem mínima de tipos)
	if err := s.validatePayloadAgainstBenchmark(ctx, benchmarkID, input.Payload); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	item, err := s.repository.Create(ctx, orgID, benchmarkID, input.Name, input.Payload)
	if err != nil {
		// Tratamento específico para violação de unique constraint
		if isUniqueViolationError(err) {
			return nil, ErrDuplicateName
		}
		return nil, fmt.Errorf("storage error: %w", err)
	}

	return s.mapItemToResponse(item), nil
}

// List retorna todas as configurações de benchmark (versão leve sem payload)
func (s *Service) List(ctx context.Context, orgID, benchmarkID string) (*dto.ConfigListAllResponseDTO, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, ErrInvalidInput
	}

	items, err := s.repository.List(ctx, orgID, benchmarkID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	var data []dto.ConfigListResponseDTO
	for _, item := range items {
		data = append(data, dto.ConfigListResponseDTO{
			ID:          fmt.Sprintf("%d", item.ID),
			Name:        item.Name,
			OrgID:       item.OrgID,
			BenchmarkID: item.BenchmarkID,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &dto.ConfigListAllResponseDTO{Data: data}, nil
}

// GetByID retorna todas as configurações de um benchmark específico
func (s *Service) GetByID(ctx context.Context, orgID, benchmarkID, id string) (*dto.ConfigResponseDTO, error) {
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

func (s *Service) Update(ctx context.Context, orgID, benchmarkID, id string, input *dto.ConfigUpdateRequestDTO) (*dto.ConfigResponseDTO, error) {
	if orgID == "" || benchmarkID == "" || id == "" {
		return nil, ErrInvalidInput
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

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

	// Atualizar payload e name
	item, err := s.repository.Update(ctx, intID, input.Name, input.Payload)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil, ErrDuplicateName
		}
		return nil, fmt.Errorf("update error: %w", err)
	}

	return s.mapItemToResponse(item), nil
}

func (s *Service) Delete(ctx context.Context, orgID, benchmarkID, id string) error {
	if orgID == "" || benchmarkID == "" || id == "" {
		return ErrInvalidInput
	}

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

// mapItemToResponse converte VarConfig para VarConfigResponse
func (s *Service) mapItemToResponse(item *models.VarConfig) *dto.ConfigResponseDTO {
	var payload map[string]any
	if item.Payload != nil {
		if err := json.Unmarshal(item.Payload, &payload); err != nil {
			payload = nil
		}
	}

	return &dto.ConfigResponseDTO{
		ID:          fmt.Sprintf("%d", item.ID),
		Name:        item.Name,
		OrgID:       item.OrgID,
		BenchmarkID: item.BenchmarkID,
		Payload:     payload,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	}
}

// payloadValidator encapsula os dados que precisam ser validados
// contra um esquema dinâmico. Implementa restwrapper.Validatable.
type payloadValidator struct {
	payload   map[string]any
	schemaDef map[string]dtoBench.SchemaFieldDefinition
}

func (v *payloadValidator) Validate() error {
	// Validar cada campo requerido e segundo sua definição
	for fieldName, def := range v.schemaDef {
		val, exists := v.payload[fieldName]
		if def.Required && !exists {
			return fmt.Errorf("field '%s' is required: %w", fieldName, ErrPayloadValidation)
		}
		if exists {
			// Usar o método Validate() da definição de campo
			if err := def.Validate(fieldName, val); err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}
	return nil
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

	validator := &payloadValidator{
		payload:   payload,
		schemaDef: schemaDef,
	}
	
	if err := validator.Validate(); err != nil {
		return err
	}

	return nil
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