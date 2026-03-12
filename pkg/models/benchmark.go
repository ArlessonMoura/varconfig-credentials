package models

import (
	"encoding/json"
	"time"
)

// BenchmarkSchema representa o schema no banco relacional PostgreSQL
type BenchmarkSchema struct {
  ID               int64            `gorm:"primaryKey;autoIncrement" json:"id"`
  Name             string           `gorm:"type:varchar(255);not null;unique" json:"name"`
  Version          string           `gorm:"type:varchar(50);not null" json:"version"`
  Schema json.RawMessage						`gorm:"type:jsonb;column:schema_definition" json:"schema_definition"`
  CreatedAt        time.Time        `gorm:"autoCreateTime" json:"created_at"`
  UpdatedAt        time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BenchmarkSchema) TableName() string {
	return "benchmark_schemas"
}

// BenchmarkSchemaCreateRequest representa os dados para criação de um novo schema
type BenchmarkSchemaCreateRequest struct {
	Name       string            `json:"name"`
	SchemaBody map[string]string `json:"schema_body"`
}

// // BenchmarkSchemaUpdateRequest representa os dados para atualização de um schema existente
// type BenchmarkSchemaUpdateRequest struct {
// 	ID         string            `json:"id"`
// 	Name       string            `json:"name"`
// 	SchemaBody map[string]string `json:"schema_body"`
// }

