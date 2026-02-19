package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto "projeto-crud-credentials/dto/config"
	ports "projeto-crud-credentials/internal/service"
	models "projeto-crud-credentials/pkg/models/config"

	"github.com/google/uuid"
)

// package config provides domain services for managing variable configurations.
// Erros de domínio básicos
var (
	ErrNotFound     = errors.New("varconfig not found")
	ErrInvalidInput = errors.New("invalid input: orgID and benchmarkID are required")
)

type Service struct {
	repository ports.IVarConfigRepository
}

func NewService(repository ports.IVarConfigRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// Create cria um novo VarConfig gerando PK e SK com ID único (UUID v4)
func (s *Service) Create(ctx context.Context, orgID string, benchmarkID string, input *dto.CreateVarConfigRequest) (*dto.VarConfigResponse, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, ErrInvalidInput
	}

	// ID único usando UUID v4
	id := uuid.NewString()
	now := time.Now().Format(time.RFC3339)

	item := models.VarConfigItem{
		PK:          fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID),
		SK:          fmt.Sprintf("VARCONFIG#%s", id),
		ID:          id,
		OrgID:       orgID,
		BenchmarkID: benchmarkID,
		Payload:     input.Payload,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repository.Create(ctx, &item); err != nil {
		return nil, fmt.Errorf("storage error: %w", err)
	}

	return s.mapItemToResponse(&item), nil
}

// List retorna todas as configurações de um benchmark
func (s *Service) List(ctx context.Context, orgID, benchmarkID string) (*dto.ListVarConfigsResponse, error) {
	pk := fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID)

	items, err := s.repository.List(ctx, pk)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	var data []dto.VarConfigResponse
	for _, item := range items {
		data = append(data, *s.mapListItemToResponse(item))
	}

	return &dto.ListVarConfigsResponse{Data: data}, nil
}

// GetByID busca um item específico usando PK e SK
func (s *Service) GetByID(ctx context.Context, orgID, benchmarkID, id string) (*dto.VarConfigResponse, error) {
	pk := fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID)
	sk := fmt.Sprintf("VARCONFIG#%s", id)

	item, err := s.repository.GetByID(ctx, pk, sk)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if item == nil {
		return nil, ErrNotFound
	}

	return s.mapItemToResponse(item), nil
}

// Update atualiza o payload
func (s *Service) Update(ctx context.Context, orgID, benchmarkID, id string, input *dto.UpdateVarConfigRequest) (*dto.VarConfigResponse, error) {
	pk := fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID)
	sk := fmt.Sprintf("VARCONFIG#%s", id)

	existing, err := s.repository.GetByID(ctx, pk, sk)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	existing.Payload = input.Payload
	existing.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := s.repository.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update error: %w", err)
	}

	return s.mapItemToResponse(existing), nil
}

// Delete remove via chaves compostas
func (s *Service) Delete(ctx context.Context, orgID, benchmarkID, id string) error {
	pk := fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID)
	sk := fmt.Sprintf("VARCONFIG#%s", id)

	return s.repository.Delete(ctx, pk, sk)
}

func (s *Service) mapItemToResponse(item *models.VarConfigItem) *dto.VarConfigResponse {
	return &dto.VarConfigResponse{
		ID:          item.ID,
		OrgID:       item.OrgID,
		BenchmarkID: item.BenchmarkID,
		Payload:     item.Payload,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func (s *Service) mapListItemToResponse(item *models.VarConfigListItem) *dto.VarConfigResponse {
	return &dto.VarConfigResponse{
		ID:          item.ID,
		OrgID:       item.OrgID,
		BenchmarkID: item.BenchmarkID,
		Payload:     nil, // Payload não disponível na listagem
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
