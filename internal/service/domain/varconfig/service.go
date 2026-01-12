package varconfig

import (
	"fmt"
	"time"
)

// Service implementa a lógica de negócio para gerenciar VarConfigs
type Service struct {
	repository VarConfigRepository
}

// NewService cria uma nova instância do service de VarConfig
func NewService(repository VarConfigRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// Create cria um novo VarConfig
func (s *Service) Create(config VarConfig) (VarConfig, error) {
	if config.OrgID <= 0 {
		return VarConfig{}, fmt.Errorf("orgID deve ser maior que zero")
	}

	if config.BenchmarkID == "" {
		return VarConfig{}, fmt.Errorf("benchmarkID não pode estar vazio")
	}

	if config.Payload == nil {
		return VarConfig{}, fmt.Errorf("payload não pode ser nulo")
	}

	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	return s.repository.Save(config)
}

// GetByID obtém um VarConfig específico
func (s *Service) GetByID(orgID int64, benchmarkID string, id int64) (VarConfig, error) {
	if orgID <= 0 {
		return VarConfig{}, fmt.Errorf("orgID deve ser maior que zero")
	}

	if benchmarkID == "" {
		return VarConfig{}, fmt.Errorf("benchmarkID não pode estar vazio")
	}

	if id <= 0 {
		return VarConfig{}, fmt.Errorf("id deve ser maior que zero")
	}

	return s.repository.FindByID(orgID, benchmarkID, id)
}

// ListByBenchmark lista todos os VarConfigs de um benchmark
func (s *Service) ListByBenchmark(orgID int64, benchmarkID string) ([]VarConfig, error) {
	if orgID <= 0 {
		return nil, fmt.Errorf("orgID deve ser maior que zero")
	}

	if benchmarkID == "" {
		return nil, fmt.Errorf("benchmarkID não pode estar vazio")
	}

	return s.repository.FindAllByBenchmark(orgID, benchmarkID)
}

// Update atualiza um VarConfig existente
func (s *Service) Update(config VarConfig) (VarConfig, error) {
	if config.ID <= 0 {
		return VarConfig{}, fmt.Errorf("id deve ser maior que zero")
	}

	if config.OrgID <= 0 {
		return VarConfig{}, fmt.Errorf("orgID deve ser maior que zero")
	}

	if config.BenchmarkID == "" {
		return VarConfig{}, fmt.Errorf("benchmarkID não pode estar vazio")
	}

	if config.Payload == nil {
		return VarConfig{}, fmt.Errorf("payload não pode ser nulo")
	}

	config.UpdatedAt = time.Now()

	return s.repository.Update(config)
}

// Delete remove um VarConfig
func (s *Service) Delete(orgID int64, benchmarkID string, id int64) error {
	if orgID <= 0 {
		return fmt.Errorf("orgID deve ser maior que zero")
	}

	if benchmarkID == "" {
		return fmt.Errorf("benchmarkID não pode estar vazio")
	}

	if id <= 0 {
		return fmt.Errorf("id deve ser maior que zero")
	}

	return s.repository.Delete(orgID, benchmarkID, id)
}
