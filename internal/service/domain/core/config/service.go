package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dto "projeto-crud-credentials/dto/config"
	ports "projeto-crud-credentials/internal/service"
	models "projeto-crud-credentials/pkg/models/config"
)

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

// Create cria um novo VarConfig com ID gerado automaticamente pelo PostgreSQL
func (s *Service) Create(ctx context.Context, orgID string, benchmarkID string, input *dto.CreateVarConfigRequest) (*dto.VarConfigResponse, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, ErrInvalidInput
	}

	item, err := s.repository.Create(ctx, orgID, benchmarkID, input.Payload)
	if err != nil {
		return nil, fmt.Errorf("storage error: %w", err)
	}

	return s.mapItemToResponse(item), nil
}

// List retorna todas as configurações de um benchmark específico
func (s *Service) List(ctx context.Context, orgID, benchmarkID string) (*dto.ListVarConfigsResponse, error) {
	if orgID == "" || benchmarkID == "" {
		return nil, ErrInvalidInput
	}

	items, err := s.repository.List(ctx, orgID, benchmarkID)
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	var data []dto.VarConfigResponse
	for _, item := range items {
		data = append(data, *s.mapItemToResponse(item))
	}

	return &dto.ListVarConfigsResponse{Data: data}, nil
}

// GetByID busca um item específico pelo ID
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
