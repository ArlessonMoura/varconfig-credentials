package varconfig

// VarConfigRepository define as operações de persistência de VarConfig
type VarConfigRepository interface {
FindAllByBenchmark(orgID int64, benchmarkID string) ([]VarConfig, error)
Save(config VarConfig) (VarConfig, error)
FindByID(orgID int64, benchmarkID string, id int64) (VarConfig, error)
Update(config VarConfig) (VarConfig, error)
Delete(orgID int64, benchmarkID string, id int64) error
}

// VarConfigUseCase define as operações de negócio para VarConfig
type VarConfigUseCase interface {
Create(config VarConfig) (VarConfig, error)
GetByID(orgID int64, benchmarkID string, id int64) (VarConfig, error)
ListByBenchmark(orgID int64, benchmarkID string) ([]VarConfig, error)
Update(config VarConfig) (VarConfig, error)
Delete(orgID int64, benchmarkID string, id int64) error
}
