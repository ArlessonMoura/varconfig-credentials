package benchmarkschemas

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	models "projeto-crud-credentials/pkg/models/benchmarkschema"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, schema *models.BenchmarkSchemaRelational) error {
	query := `
		INSERT INTO benchmark_schemas (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	
	// Executa a query e já mapeia o ID gerado pelo Serial/Auto-increment de volta para a struct
	err := r.db.QueryRowContext(ctx, query, schema.Name, now, now).Scan(
		&schema.ID, 
		&schema.CreatedAt, 
		&schema.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert benchmark schema: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id *int64) error {
	query := `DELETE FROM benchmark_schemas WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, *id)
	if err != nil {
		return fmt.Errorf("failed to delete benchmark schema for rollback: %w", err)
	}

	return nil
}