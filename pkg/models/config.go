package models

import (
	"encoding/json"
	"time"
)

// VarConfig representa a configuração de variáveis no banco PostgreSQL
type VarConfig struct {
  ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
  Name        string          `gorm:"type:varchar(255);not null;uniqueIndex:idx_varcfg_org_bench_name" json:"name"` 
  OrgID       string          `gorm:"not null;uniqueIndex:idx_varcfg_org_bench_name" json:"org_id"`
  BenchmarkID string          `gorm:"not null;uniqueIndex:idx_varcfg_org_bench_name;constraint:OnDelete:CASCADE" json:"benchmark_id"`
  Benchmark   *BenchmarkSchema `gorm:"foreignKey:BenchmarkID;constraint:OnDelete:CASCADE" json:"benchmark,omitempty"`
  Payload     json.RawMessage `gorm:"type:jsonb;column:payload" json:"payload"`
  CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
  UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (VarConfig) TableName() string {
	return "var_configs"
}

// VarConfigCreateRequest representa os dados para criação de uma nova configuração
type VarConfigCreateRequest struct {
	OrgID       string         `json:"org_id"`
	BenchmarkID string         `json:"benchmark_id"`
	Payload     map[string]any `json:"payload"`
}

// VarConfigUpdateRequest representa os dados para atualização de uma configuração existente
type VarConfigUpdateRequest struct {
	ID          string         `json:"id"`
	OrgID       string         `json:"org_id"`
	BenchmarkID string         `json:"benchmark_id"`
	Payload     map[string]any `json:"payload"`
}
