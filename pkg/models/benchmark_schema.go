package models

import "time"

type BenchmarkSchemaRelational struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BenchmarkSchemaRelational) TableName() string {
	return "benchmark_schemas"
}

type BenchmarkSchemaNoSQL struct {
	PK         string         `dynamodbav:"PK"` 
  SK         string         `dynamodbav:"SK"`
	ID        string         `dynamodbav:"ID"`        // O ID vindo do Relacional (como string)
	SchemaBody map[string]string `dynamodbav:"schema_body"`
	CreatedAt string         `dynamodbav:"created_at"`
}