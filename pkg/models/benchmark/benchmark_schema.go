package models

import "time"

// BenchmarkSchemaPostgreSQL representa o schema no banco relacional PostgreSQL
type BenchmarkSchemaPostgreSQL struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
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

// BenchmarkSchemaDynamoDB representa o schema completo no DynamoDB
type BenchmarkSchemaDynamoDB struct {
	PK         string            `dynamodbav:"PK"` // SCHEMA#<id>
	ID         string            `dynamodbav:"ID"` // O ID vindo do Relacional (como string)
	SchemaBody map[string]string `dynamodbav:"schema_body"`
	CreatedAt  string            `dynamodbav:"created_at"`
}

// BenchmarkSchemaListDynamoDB representa os campos essenciais para listagem no DynamoDB (sem o schema_body pesado)
type BenchmarkSchemaListDynamoDB struct {
	ID        string `dynamodbav:"ID"`
	CreatedAt string `dynamodbav:"created_at"`
}