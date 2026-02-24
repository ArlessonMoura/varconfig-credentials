package models

import (
	"encoding/json"
	"time"
)

// BenchmarkSchemaPostgreSQL representa o schema no banco relacional PostgreSQL
type BenchmarkSchemaPostgreSQL struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Schema    json.RawMessage `gorm:"type:jsonb;column:schema_body" json:"schema_body"`
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BenchmarkSchemaPostgreSQL) TableName() string {
	return "benchmark_schemas"
}

// BenchmarkSchemaCreateRequest representa os dados para criação de um novo schema
type BenchmarkSchemaCreateRequest struct {
	Name      string `json:"name"`
	SchemaBody map[string]string `json:"schema_body"`
}

// // BenchmarkSchemaUpdateRequest representa os dados para atualização de um schema existente
// type BenchmarkSchemaUpdateRequest struct {
// 	ID         string            `json:"id"`
// 	Name       string            `json:"name"`
// 	SchemaBody map[string]string `json:"schema_body"`
// }

// DynamoDB artifacts removed — benchmark schemas are stored in Postgres (jsonb).