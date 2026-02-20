package service

import (
	"context"
	modelBenchmark "projeto-crud-credentials/pkg/models/benchmark"
	modelConfig "projeto-crud-credentials/pkg/models/config"
)


type IVarConfigRepository interface {
	Create(ctx context.Context, orgID string, benchmarkID string, payload map[string]any) (*modelConfig.VarConfigPostgreSQL, error)

	List(ctx context.Context, orgID string, benchmarkID string) ([]*modelConfig.VarConfigPostgreSQL, error)

	GetByID(ctx context.Context, id int64) (*modelConfig.VarConfigPostgreSQL, error)

	Update(ctx context.Context, id int64, payload map[string]any) (*modelConfig.VarConfigPostgreSQL, error)

	Delete(ctx context.Context, id int64) error
}


type IBenchmarkRepository interface {
	Create(ctx context.Context, schema *modelBenchmark.BenchmarkSchemaPostgreSQL) error
	
	List(ctx context.Context) ([]*modelBenchmark.BenchmarkSchemaPostgreSQL, error)
	
	GetByID(ctx context.Context, id int64) (*modelBenchmark.BenchmarkSchemaPostgreSQL, error)
	
	Delete(ctx context.Context, id *int64) error
}

