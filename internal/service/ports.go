package service

import (
	"context"
	"projeto-crud-credentials/pkg/models"
)

type IVarConfigRepository interface {
	Create(ctx context.Context, orgID string, benchmarkID string, payload map[string]any) (*models.VarConfig, error)

	List(ctx context.Context, orgID string, benchmarkID string) ([]*models.VarConfig, error)

	GetByID(ctx context.Context, id int64) (*models.VarConfig, error)

	Update(ctx context.Context, id int64, payload map[string]any) (*models.VarConfig, error)

	Delete(ctx context.Context, id int64) error
}

type IBenchmarkRepository interface {
	Create(ctx context.Context, schema *models.BenchmarkSchema) error

	List(ctx context.Context) ([]*models.BenchmarkSchema, error)

	GetByID(ctx context.Context, id int64) (*models.BenchmarkSchema, error)

	Delete(ctx context.Context, id *int64) error
}
