package benchmark

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"projeto-crud-credentials/pkg/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create cria o registro usando uma transação GORM
func (r *Repository) Create(ctx context.Context, schema *models.BenchmarkSchema) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(schema).Error; err != nil {
			return fmt.Errorf("failed to insert benchmark schema: %w", err)
		}
		return nil
	})
}

func (r *Repository) List(ctx context.Context) ([]*models.BenchmarkSchema, error) {
	var items []*models.BenchmarkSchema

	// Selecionar apenas os campos de metadados para evitar transferência do JSONB (schema_definition)
	if err := r.db.WithContext(ctx).
		Select("id", "name", "version", "created_at", "updated_at").
		Order("created_at desc").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list benchmark schemas: %w", err)
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*models.BenchmarkSchema, error) {
	var item models.BenchmarkSchema
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get benchmark schema: %w", err)
	}
	return &item, nil
}

func (r *Repository) Delete(ctx context.Context, id *int64) error {
	if err := r.db.WithContext(ctx).Delete(&models.BenchmarkSchema{}, *id).Error; err != nil {
		return fmt.Errorf("failed to delete benchmark schema: %w", err)
	}
	return nil
}
