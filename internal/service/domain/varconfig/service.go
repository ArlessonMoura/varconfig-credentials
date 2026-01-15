// Package varconfig implements business logic for VarConfig management.
//
// This service layer contains all business rules and validation logic
// for VarConfig operations, following clean architecture principles.
package varconfig

import (
	"context"
	"time"

	"projeto-crud-credencials/internal/common/errors"
	"projeto-crud-credencials/internal/common/logger"
	"projeto-crud-credencials/internal/common/metrics"
)

// Service implementa a lógica de negócio para gerenciar VarConfigs
type Service struct {
	repository VarConfigRepository
	logger     logger.Logger
	metrics    metrics.Collector
}

// NewService cria uma nova instância do service de VarConfig
func NewService(repository VarConfigRepository, logger logger.Logger, metrics metrics.Collector) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
		metrics:    metrics,
	}
}

// Create cria um novo VarConfig
func (s *Service) Create(ctx context.Context, config VarConfig) (VarConfig, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordLatency(ctx, metrics.ServiceVarConfigCreateLatency, time.Since(start))
	}()

	if err := s.validateCreateInput(config); err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigCreateFailures)
		return VarConfig{}, err
	}

	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	s.logger.Info(ctx, "Criando VarConfig para orgID=%d, benchmarkID=%s", config.OrgID, config.BenchmarkID)

	result, err := s.repository.Save(ctx, config)
	if err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigCreateFailures)
		s.logger.Error(ctx, "Erro ao criar VarConfig: %v", err)
		return VarConfig{}, errors.NewStorageError("Erro ao criar VarConfig", err)
	}

	s.metrics.RecordSuccess(ctx, metrics.ServiceVarConfigCreateLatency)
	s.logger.Info(ctx, "VarConfig criado com sucesso, ID=%d", result.ID)
	return result, nil
}

// GetByID obtém um VarConfig específico
func (s *Service) GetByID(ctx context.Context, orgID int64, benchmarkID string, id int64) (VarConfig, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordLatency(ctx, metrics.ServiceVarConfigGetLatency, time.Since(start))
	}()

	if err := s.validateGetInput(orgID, benchmarkID, id); err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigGetFailures)
		return VarConfig{}, err
	}

	s.logger.Debug(ctx, "Buscando VarConfig orgID=%d, benchmarkID=%s, id=%d", orgID, benchmarkID, id)

	result, err := s.repository.FindByID(ctx, orgID, benchmarkID, id)
	if err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigGetFailures)
		s.logger.Error(ctx, "Erro ao buscar VarConfig: %v", err)
		return VarConfig{}, errors.NewStorageError("Erro ao buscar VarConfig", err)
	}

	if result.ID == 0 {
		return VarConfig{}, errors.NewNotFoundError("VarConfig")
	}

	s.metrics.RecordSuccess(ctx, metrics.ServiceVarConfigGetLatency)
	s.logger.Debug(ctx, "VarConfig encontrado com sucesso")
	return result, nil
}

// ListByBenchmark lista todos os VarConfigs de um benchmark
func (s *Service) ListByBenchmark(ctx context.Context, orgID int64, benchmarkID string) ([]VarConfig, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordLatency(ctx, metrics.ServiceVarConfigListLatency, time.Since(start))
	}()

	if err := s.validateListInput(orgID, benchmarkID); err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigListFailures)
		return nil, err
	}

	s.logger.Debug(ctx, "Listando VarConfigs para orgID=%d, benchmarkID=%s", orgID, benchmarkID)

	results, err := s.repository.FindAllByBenchmark(ctx, orgID, benchmarkID)
	if err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigListFailures)
		s.logger.Error(ctx, "Erro ao listar VarConfigs: %v", err)
		return nil, errors.NewStorageError("Erro ao listar VarConfigs", err)
	}

	s.metrics.RecordSuccess(ctx, metrics.ServiceVarConfigListLatency)
	s.logger.Info(ctx, "Listados %d VarConfigs com sucesso", len(results))
	return results, nil
}

// Update atualiza um VarConfig existente
func (s *Service) Update(ctx context.Context, config VarConfig) (VarConfig, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordLatency(ctx, metrics.ServiceVarConfigUpdateLatency, time.Since(start))
	}()

	if err := s.validateUpdateInput(config); err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigUpdateFailures)
		return VarConfig{}, err
	}

	// Verificar se o VarConfig existe antes de atualizar
	existing, err := s.repository.FindByID(ctx, config.OrgID, config.BenchmarkID, config.ID)
	if err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigUpdateFailures)
		return VarConfig{}, errors.NewStorageError("Erro ao verificar VarConfig existente", err)
	}

	if existing.ID == 0 {
		return VarConfig{}, errors.NewNotFoundError("VarConfig")
	}

	config.UpdatedAt = time.Now()

	s.logger.Info(ctx, "Atualizando VarConfig ID=%d para orgID=%d, benchmarkID=%s", config.ID, config.OrgID, config.BenchmarkID)

	result, err := s.repository.Update(ctx, config)
	if err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigUpdateFailures)
		s.logger.Error(ctx, "Erro ao atualizar VarConfig: %v", err)
		return VarConfig{}, errors.NewStorageError("Erro ao atualizar VarConfig", err)
	}

	s.metrics.RecordSuccess(ctx, metrics.ServiceVarConfigUpdateLatency)
	s.logger.Info(ctx, "VarConfig atualizado com sucesso")
	return result, nil
}

// Delete remove um VarConfig
func (s *Service) Delete(ctx context.Context, orgID int64, benchmarkID string, id int64) error {
	start := time.Now()
	defer func() {
		s.metrics.RecordLatency(ctx, metrics.ServiceVarConfigDeleteLatency, time.Since(start))
	}()

	if err := s.validateDeleteInput(orgID, benchmarkID, id); err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigDeleteFailures)
		return err
	}

	// Verificar se o VarConfig existe antes de deletar
	existing, err := s.repository.FindByID(ctx, orgID, benchmarkID, id)
	if err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigDeleteFailures)
		return errors.NewStorageError("Erro ao verificar VarConfig existente", err)
	}

	if existing.ID == 0 {
		return errors.NewNotFoundError("VarConfig")
	}

	s.logger.Info(ctx, "Deletando VarConfig ID=%d para orgID=%d, benchmarkID=%s", id, orgID, benchmarkID)

	err = s.repository.Delete(ctx, orgID, benchmarkID, id)
	if err != nil {
		s.metrics.RecordFailure(ctx, metrics.StorageVarConfigDeleteFailures)
		s.logger.Error(ctx, "Erro ao deletar VarConfig: %v", err)
		return errors.NewStorageError("Erro ao deletar VarConfig", err)
	}

	s.metrics.RecordSuccess(ctx, metrics.ServiceVarConfigDeleteLatency)
	s.logger.Info(ctx, "VarConfig deletado com sucesso")
	return nil
}

// Métodos de validação privados

func (s *Service) validateCreateInput(config VarConfig) error {
	if config.OrgID <= 0 {
		return errors.NewValidationError("orgID deve ser maior que zero", map[string]any{"field": "orgID", "value": config.OrgID})
	}
	if config.BenchmarkID == "" {
		return errors.NewValidationError("benchmarkID não pode estar vazio", map[string]any{"field": "benchmarkID"})
	}
	if config.Payload == nil {
		return errors.NewValidationError("payload não pode ser nulo", map[string]any{"field": "payload"})
	}
	return nil
}

func (s *Service) validateGetInput(orgID int64, benchmarkID string, id int64) error {
	if orgID <= 0 {
		return errors.NewValidationError("orgID deve ser maior que zero", map[string]any{"field": "orgID", "value": orgID})
	}
	if benchmarkID == "" {
		return errors.NewValidationError("benchmarkID não pode estar vazio", map[string]any{"field": "benchmarkID"})
	}
	if id <= 0 {
		return errors.NewValidationError("id deve ser maior que zero", map[string]any{"field": "id", "value": id})
	}
	return nil
}

func (s *Service) validateListInput(orgID int64, benchmarkID string) error {
	if orgID <= 0 {
		return errors.NewValidationError("orgID deve ser maior que zero", map[string]any{"field": "orgID", "value": orgID})
	}
	if benchmarkID == "" {
		return errors.NewValidationError("benchmarkID não pode estar vazio", map[string]any{"field": "benchmarkID"})
	}
	return nil
}

func (s *Service) validateUpdateInput(config VarConfig) error {
	if config.ID <= 0 {
		return errors.NewValidationError("id deve ser maior que zero", map[string]any{"field": "id", "value": config.ID})
	}
	return s.validateCreateInput(config)
}

func (s *Service) validateDeleteInput(orgID int64, benchmarkID string, id int64) error {
	return s.validateGetInput(orgID, benchmarkID, id)
}
