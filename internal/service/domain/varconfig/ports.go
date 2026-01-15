package varconfig

import "context"

// VarConfigRepository define as operações de persistência de VarConfig
type VarConfigRepository interface {
	FindAllByBenchmark(ctx context.Context, orgID int64, benchmarkID string) ([]VarConfig, error)
	Save(ctx context.Context, config VarConfig) (VarConfig, error)
	FindByID(ctx context.Context, orgID int64, benchmarkID string, id int64) (VarConfig, error)
	Update(ctx context.Context, config VarConfig) (VarConfig, error)
	Delete(ctx context.Context, orgID int64, benchmarkID string, id int64) error
}

// VarConfigUseCase define as operações de negócio para VarConfig
type VarConfigUseCase interface {
	Create(ctx context.Context, config VarConfig) (VarConfig, error)
	GetByID(ctx context.Context, orgID int64, benchmarkID string, id int64) (VarConfig, error)
	ListByBenchmark(ctx context.Context, orgID int64, benchmarkID string) ([]VarConfig, error)
	Update(ctx context.Context, config VarConfig) (VarConfig, error)
	Delete(ctx context.Context, orgID int64, benchmarkID string, id int64) error
}
