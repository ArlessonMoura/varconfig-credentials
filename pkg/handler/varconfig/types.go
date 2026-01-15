package varconfig

import (
	"context"
	svcvarconfig "projeto-crud-credencials/internal/service/domain/varconfig"
)

// VarConfigUseCase define as operações que o handler espera do service
type VarConfigUseCase interface {
	Create(ctx context.Context, config svcvarconfig.VarConfig) (svcvarconfig.VarConfig, error)
	GetByID(ctx context.Context, orgID int64, benchmarkID string, id int64) (svcvarconfig.VarConfig, error)
	ListByBenchmark(ctx context.Context, orgID int64, benchmarkID string) ([]svcvarconfig.VarConfig, error)
	Update(ctx context.Context, config svcvarconfig.VarConfig) (svcvarconfig.VarConfig, error)
	Delete(ctx context.Context, orgID int64, benchmarkID string, id int64) error
}
