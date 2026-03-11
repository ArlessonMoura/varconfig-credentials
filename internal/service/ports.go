package service

import (
	"context"
	"projeto-crud-credentials/pkg/models"
)

type IVarConfigRepository interface {
	Create(ctx context.Context, orgID string, benchmarkID string, payload map[string]any) (*models.VarConfigPostgreSQL, error)

	List(ctx context.Context, orgID string, benchmarkID string) ([]*models.VarConfigPostgreSQL, error)

	GetByID(ctx context.Context, id int64) (*models.VarConfigPostgreSQL, error)

	Update(ctx context.Context, id int64, payload map[string]any) (*models.VarConfigPostgreSQL, error)

	Delete(ctx context.Context, id int64) error
}

type IBenchmarkRepository interface {
	Create(ctx context.Context, schema *models.BenchmarkSchemaPostgreSQL) error

	List(ctx context.Context) ([]*models.BenchmarkSchemaPostgreSQL, error)

	GetByID(ctx context.Context, id int64) (*models.BenchmarkSchemaPostgreSQL, error)

	Delete(ctx context.Context, id *int64) error
}
