package models

import "time"

type BenchmarkSchemaRelational struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type BenchmarkSchemaNoSQL struct {
	ID        string         `dynamodbav:"ID"`        // O ID vindo do Relacional (como string)
	SchemaBody map[string]any `dynamodbav:"schema_body"`
	CreatedAt string         `dynamodbav:"created_at"`
}