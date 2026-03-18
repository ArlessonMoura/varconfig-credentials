package benchmark

import (
	"context"
	"encoding/json"
	"fmt"

	dtoBenchmark "projeto-crud-credentials/dto/benchmark"
	"projeto-crud-credentials/internal/common"
	"projeto-crud-credentials/internal/common/helpers"
	ports "projeto-crud-credentials/internal/service"
	svcbenchmark "projeto-crud-credentials/pkg/handler"
	"projeto-crud-credentials/pkg/models"
)

type Service struct {
	relationalRepo ports.IBenchmarkRepository
}

func NewService(relationalRepo ports.IBenchmarkRepository) *Service {
	return &Service{relationalRepo: relationalRepo}
}

func (s *Service) Create(ctx context.Context, schemaRequest *dtoBenchmark.BenchmarkCreateRequestDTO) (*models.BenchmarkSchema, error) {
	
	if err := schemaRequest.Validate(); err != nil {
		return nil, err
	}

	if len(schemaRequest.Schema) == 0 {
		return nil, common.ErrBenchmarkEmptySchema
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
		return nil, helpers.ValidateUniqueViolationError(err, common.ErrBenchmarkDuplicateName)
	}

	return schema, nil
}

// List recupera todos os schemas disponíveis 
func (s *Service) List(ctx context.Context) (*dtoBenchmark.BenchmarkListResponseDTO, error) {
	items, err := s.relationalRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas from repository: %w", err)
	}

	if items == nil {
		return &dtoBenchmark.BenchmarkListResponseDTO{Data: []dtoBenchmark.BenchmarkResponseDTO{}, Count: 0}, nil
	}

	var data []dtoBenchmark.BenchmarkResponseDTO
	for _, item := range items {
		data = append(data, *helpers.MapBenchmarkSchemaToResponse(item))
	}

	return &dtoBenchmark.BenchmarkListResponseDTO{
		Data:  data,
		Count: len(data),
	}, nil
}

// GetByID busca um schema pelo ID
func (s *Service) GetByID(ctx context.Context, id string) (*dtoBenchmark.BenchmarkResponseDTO, error) {
		var intID int64
	if _, err := fmt.Sscan(id, &intID); err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	item, err := s.relationalRepo.GetByID(ctx, intID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema from repository: %w", err)
	}
	if item == nil {
		return nil, common.ErrBenchmarkSchemaNotFound
	}

	return helpers.MapBenchmarkSchemaToResponseWithSchema(item), nil
}

// Delete remove um benchmark schema e todas as configurações associadas 
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

		if err := s.relationalRepo.Delete(ctx, &intID); err != nil {
		return fmt.Errorf("failed to delete benchmark schema: %w", err)
	}

	return nil
}

var _ svcbenchmark.IBenchmarkSchemaService = (*Service)(nil)