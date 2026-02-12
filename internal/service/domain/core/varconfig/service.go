// internal/service/domain/core/varconfig/service.go - CORRIGIDO
package varconfig

import (
	"context"
	"errors"
	"fmt"
	repository "projeto-crud-credencials/internal/service"
	"time"

	dto "projeto-crud-credencials/dto/varconfig_dto"
	models "projeto-crud-credencials/pkg/models/varconfig"
)

// Erros de domínio básicos
var (
	ErrNotFound     = errors.New("varconfig not found")
	ErrInvalidInput = errors.New("invalid input: orgID and benchmarkID are required")
)

type Service struct {
	repository repository.IVarConfigRepository
}

func NewService(repository repository.IVarConfigRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// Create cria um novo VarConfig gerando PK e SK com ID em nanosegundos
func (s *Service) Create(ctx context.Context, orgID string, benchmarkID string, input dto.CreateVarConfigRequest) (dto.VarConfigResponse, error) {
	if orgID == "" || benchmarkID == "" {
		return dto.VarConfigResponse{}, ErrInvalidInput
	}

	// ID em nanosegundos (Unix nano timestamp)
	id := fmt.Sprintf("%d", time.Now().UnixNano())
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

	if err := s.repository.Save(ctx, item); err != nil {
		return dto.VarConfigResponse{}, fmt.Errorf("storage error: %w", err)
	}

	return s.mapItemToResponse(item), nil
}

// GetByID busca um item específico usando PK e SK
func (s *Service) GetByID(ctx context.Context, orgID, benchmarkID, id string) (dto.VarConfigResponse, error) {
	pk := fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID)
	sk := fmt.Sprintf("VARCONFIG#%s", id)

	item, err := s.repository.Get(ctx, pk, sk)
	if err != nil {
		return dto.VarConfigResponse{}, fmt.Errorf("repository error: %w", err)
	}
	if item == nil {
		return dto.VarConfigResponse{}, ErrNotFound
	}

	return s.mapItemToResponse(*item), nil
}

// List retorna todas as configurações de um benchmark
func (s *Service) List(ctx context.Context, orgID, benchmarkID string) (dto.ListVarConfigResponse, error) {
	pk := fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID)

	items, err := s.repository.ListByPK(ctx, pk)
	if err != nil {
		return dto.ListVarConfigResponse{}, fmt.Errorf("repository error: %w", err)
	}

	var data []dto.VarConfigResponse
	for _, item := range items {
		data = append(data, s.mapItemToResponse(item))
	}

	return dto.ListVarConfigResponse{Data: data}, nil
}

// Update atualiza o payload
func (s *Service) Update(ctx context.Context, orgID, benchmarkID, id string, input dto.UpdateVarConfigRequest) (dto.VarConfigResponse, error) {
	pk := fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID)
	sk := fmt.Sprintf("VARCONFIG#%s", id)

	existing, err := s.repository.Get(ctx, pk, sk)
	if err != nil {
		return dto.VarConfigResponse{}, fmt.Errorf("repository error: %w", err)
	}
	if existing == nil {
		return dto.VarConfigResponse{}, ErrNotFound
	}

	existing.Payload = input.Payload
	existing.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := s.repository.Save(ctx, *existing); err != nil {
		return dto.VarConfigResponse{}, fmt.Errorf("update error: %w", err)
	}

	return s.mapItemToResponse(*existing), nil
}

// Delete remove via chaves compostas
func (s *Service) Delete(ctx context.Context, orgID, benchmarkID, id string) error {
	pk := fmt.Sprintf("ORG#%s#BENCH#%s", orgID, benchmarkID)
	sk := fmt.Sprintf("VARCONFIG#%s", id)

	return s.repository.Delete(ctx, pk, sk)
}

func (s *Service) mapItemToResponse(item models.VarConfigItem) dto.VarConfigResponse {
	return dto.VarConfigResponse{
		ID:          item.ID,
		OrgID:       item.OrgID,
		BenchmarkID: item.BenchmarkID,
		Payload:     item.Payload,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
