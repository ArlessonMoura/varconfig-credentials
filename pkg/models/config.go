package models

import (
	"encoding/json"
	"time"
)

// VarConfig representa a configuração de variáveis no banco PostgreSQL
type VarConfig struct {
	ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrgID       string          `gorm:"not null;index:idx_varcfg_org_bench" json:"org_id"`
	BenchmarkID string          `gorm:"not null;index:idx_varcfg_org_bench" json:"benchmark_id"`
	Payload     json.RawMessage `gorm:"type:jsonb;column:payload" json:"payload"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (VarConfig) TableName() string {
	return "var_configs"
}

// // VarConfig é a entidade de domínio e resposta da API
// type VarConfig struct {
// 	ID          string         `json:"id"`
// 	OrgID       string         `json:"org_id"`
// 	BenchmarkID string         `json:"benchmark_id"`
// 	Payload     map[string]any `json:"payload"`
// 	CreatedAt   time.Time      `json:"created_at"`
// 	UpdatedAt   time.Time      `json:"updated_at"`
// }

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
