package config

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"projeto-crud-credentials/pkg/models"
)

// Repository implementa IVarConfigRepository usando PostgreSQL + GORM
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, orgID, benchmarkID, name string, payload map[string]any) (*models.VarConfig, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	item := &models.VarConfig{
		Name:        name,
		OrgID:       orgID,
		BenchmarkID: benchmarkID,
		Payload:     payloadJSON,
	}

	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to insert var_config: %w", err)
	}

	return item, nil
}

func (r *Repository) List(ctx context.Context, orgID, benchmarkID string) ([]*models.VarConfig, error) {
	var items []*models.VarConfig

	// Selecionar apenas os campos necessários (sem payload) para reduzir transferência de dados
	if err := r.db.WithContext(ctx).
		Select("id", "name", "org_id", "benchmark_id", "created_at", "updated_at").
		Where("org_id = ? AND benchmark_id = ?", orgID, benchmarkID).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list var_configs: %w", err)
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*models.VarConfig, error) {
	var item models.VarConfig

	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get var_config: %w", err)
	}

	return &item, nil
}

func (r *Repository) Update(ctx context.Context, id int64, name string, payload map[string]any) (*models.VarConfig, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Atualizar tanto name quanto payload
	if err := r.db.WithContext(ctx).Model(&models.VarConfig{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":    name,
			"payload": payloadJSON,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update var_config: %w", err)
	}

	// Retorna o item atualizado
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&models.VarConfig{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete var_config: %w", err)
	}

	return nil
}
