package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"projeto-crud-credentials/dto/benchmark"
	"projeto-crud-credentials/dto/config"
	"projeto-crud-credentials/internal/common"
	"projeto-crud-credentials/internal/common/helpers"
	ports "projeto-crud-credentials/internal/service"
	svcvarconfig "projeto-crud-credentials/pkg/handler"
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
func (s *Service) Create(ctx context.Context, orgID string, benchmarkID string, input *config.ConfigCreateRequestDTO) (*config.ConfigResponseDTO, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, common.ErrConfigInvalidInput
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
		return nil, helpers.ValidateUniqueViolationError(err, common.ErrConfigDuplicateName)
	}

	return helpers.MapVarConfigToResponse(item), nil
}

// List retorna todas as configurações de benchmark (versão leve sem payload)
func (s *Service) List(ctx context.Context, orgID, benchmarkID string) (*config.ConfigListAllResponseDTO, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, common.ErrConfigInvalidInput
	}

	items, err := s.repository.List(ctx, orgID, benchmarkID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	var data []config.ConfigListResponseDTO
	for _, item := range items {
		data = append(data, config.ConfigListResponseDTO{
			ID:          fmt.Sprintf("%d", item.ID),
			Name:        item.Name,
			OrgID:       item.OrgID,
			BenchmarkID: fmt.Sprintf("%d", item.BenchmarkID),
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &config.ConfigListAllResponseDTO{Data: data}, nil
}

// GetByID retorna todas as configurações de um benchmark específico
func (s *Service) GetByID(ctx context.Context, orgID, benchmarkID, id string) (*config.ConfigResponseDTO, error) {
	if orgID == "" || benchmarkID == "" || id == "" {
		return nil, common.ErrConfigInvalidInput
	}

	// Converter ID string para int64
	var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// Converter benchmarkID para int64
	var benchmarkIDInt int64
	if _, err := fmt.Sscanf(benchmarkID, "%d", &benchmarkIDInt); err != nil {
		return nil, fmt.Errorf("invalid benchmark id format: %w", err)
	}

	item, err := s.repository.GetByID(ctx, intID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if item == nil {
		return nil, common.ErrConfigNotFound
	}

	// Validar que o item pertence aos IDs fornecidos
	if item.OrgID != orgID || item.BenchmarkID != benchmarkIDInt {
		return nil, common.ErrConfigNotFound
	}

	return helpers.MapVarConfigToResponse(item), nil
}

func (s *Service) Update(ctx context.Context, orgID, benchmarkID, id string, input *config.ConfigUpdateRequestDTO) (*config.ConfigResponseDTO, error) {
	if orgID == "" || benchmarkID == "" || id == "" {
		return nil, common.ErrConfigInvalidInput
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// Converter benchmarkID para int64
	var benchmarkIDInt int64
	if _, err := fmt.Sscanf(benchmarkID, "%d", &benchmarkIDInt); err != nil {
		return nil, fmt.Errorf("invalid benchmark id format: %w", err)
	}

	// Verificar que o item existe e pertence aos IDs fornecidos
	existing, err := s.repository.GetByID(ctx, intID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if existing == nil {
		return nil, common.ErrConfigNotFound
	}

	if existing.OrgID != orgID || existing.BenchmarkID != benchmarkIDInt {
		return nil, common.ErrConfigNotFound
	}

	// Validar payload contra o schema do benchmark (minima de checagem de tipos)
	if err := s.validatePayloadAgainstBenchmark(ctx, benchmarkID, input.Payload); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Atualizar payload e name
	item, err := s.repository.Update(ctx, intID, input.Name, input.Payload)
	if err != nil {
		return nil, helpers.ValidateUniqueViolationError(err, common.ErrConfigDuplicateName)
	}

	return helpers.MapVarConfigToResponse(item), nil
}

func (s *Service) Delete(ctx context.Context, orgID, benchmarkID, id string) error {
	if orgID == "" || benchmarkID == "" || id == "" {
		return common.ErrConfigInvalidInput
	}

	var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	// Converter benchmarkID para int64
	var benchmarkIDInt int64
	if _, err := fmt.Sscanf(benchmarkID, "%d", &benchmarkIDInt); err != nil {
		return fmt.Errorf("invalid benchmark id format: %w", err)
	}

	// Verificar que o item existe e pertence aos IDs fornecidos
	existing, err := s.repository.GetByID(ctx, intID)
	if err != nil {
		return fmt.Errorf("repository error: %w", err)
	}
	if existing == nil {
		return common.ErrConfigNotFound
	}

	if existing.OrgID != orgID || existing.BenchmarkID != benchmarkIDInt {
		return common.ErrConfigNotFound
	}

	return s.repository.Delete(ctx, intID)
}

// payloadValidator encapsula os dados que precisam ser validados
// contra um esquema dinâmico. Implementa restwrapper.Validatable.
type payloadValidator struct {
	payload   map[string]any
	schemaDef map[string]benchmark.SchemaFieldDefinition
}

func (v *payloadValidator) Validate() error {
	// Validar cada campo requerido e segundo sua definição
	for fieldName, def := range v.schemaDef {
		val, exists := v.payload[fieldName]
		if def.Required && !exists {
			return fmt.Errorf("campo '%s' é obrigatório: %w", fieldName, common.ErrConfigPayloadValidation)
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
		return fmt.Errorf("repository error: %w", err)
	}
	if bench == nil {
		return common.ErrBenchmarkSchemaNotFound
	}

	var schemaDef map[string]benchmark.SchemaFieldDefinition
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

var _ svcvarconfig.IVarConfigService = (*Service)(nil)
