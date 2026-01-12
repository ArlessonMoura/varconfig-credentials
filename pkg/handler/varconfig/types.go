package varconfig

import svcvarconfig "projeto-crud-credencials/internal/service/domain/varconfig"

// VarConfigUseCase define as operações que o handler espera do service
type VarConfigUseCase interface {
	Create(config svcvarconfig.VarConfig) (svcvarconfig.VarConfig, error)
	GetByID(orgID int64, benchmarkID string, id int64) (svcvarconfig.VarConfig, error)
	ListByBenchmark(orgID int64, benchmarkID string) ([]svcvarconfig.VarConfig, error)
	Update(config svcvarconfig.VarConfig) (svcvarconfig.VarConfig, error)
	Delete(orgID int64, benchmarkID string, id int64) error
}
