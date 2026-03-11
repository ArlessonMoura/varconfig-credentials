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

func (r *Repository) Create(ctx context.Context, orgID, benchmarkID string, payload map[string]any) (*models.VarConfigPostgreSQL, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	item := &models.VarConfigPostgreSQL{
		OrgID:       orgID,
		BenchmarkID: benchmarkID,
		Payload:     payloadJSON,
	}

	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to insert var_config: %w", err)
	}

	return item, nil
}

func (r *Repository) List(ctx context.Context, orgID, benchmarkID string) ([]*models.VarConfigPostgreSQL, error) {
	var items []*models.VarConfigPostgreSQL

	// Selecionar apenas os campos necessários (sem payload) para reduzir transferência de dados
	if err := r.db.WithContext(ctx).
		Select("id", "org_id", "benchmark_id", "created_at", "updated_at").
		Where("org_id = ? AND benchmark_id = ?", orgID, benchmarkID).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list var_configs: %w", err)
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*models.VarConfigPostgreSQL, error) {
	var item models.VarConfigPostgreSQL

	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get var_config: %w", err)
	}

	return &item, nil
}

func (r *Repository) Update(ctx context.Context, id int64, payload map[string]any) (*models.VarConfigPostgreSQL, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	item := &models.VarConfigPostgreSQL{Payload: payloadJSON}

	if err := r.db.WithContext(ctx).Model(&models.VarConfigPostgreSQL{}).
		Where("id = ?", id).
		Update("payload", item.Payload).Error; err != nil {
		return nil, fmt.Errorf("failed to update var_config: %w", err)
	}

	// Retorna o item atualizado
	return r.GetByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&models.VarConfigPostgreSQL{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete var_config: %w", err)
	}

	return nil
}
