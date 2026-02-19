package service

import (
	"context"
	modelBenchmark "projeto-crud-credentials/pkg/models/benchmark"
	modelConfig "projeto-crud-credentials/pkg/models/config"
)

// IVarConfigRepository define o contrato de persistência para VarConfig no PostgreSQL.
// O Service trabalha com domínio relacional (ID, orgID, benchmarkID, payload).
type IVarConfigRepository interface {
	// Create persiste um novo VarConfig e retorna o item criado com ID atribuído
	Create(ctx context.Context, orgID string, benchmarkID string, payload map[string]any) (*modelConfig.VarConfigPostgreSQL, error)

	// List retorna todas as configurações para um dado org_id e benchmark_id
	List(ctx context.Context, orgID string, benchmarkID string) ([]*modelConfig.VarConfigPostgreSQL, error)

	// GetByID busca um item específico pelo ID
	GetByID(ctx context.Context, id int64) (*modelConfig.VarConfigPostgreSQL, error)

	// Update atualiza o payload de um VarConfig
	Update(ctx context.Context, id int64, payload map[string]any) (*modelConfig.VarConfigPostgreSQL, error)

	// Delete remove um item pelo ID
	Delete(ctx context.Context, id int64) error
}

// IBenchmarkRepository define o contrato para persistência relacional (PostgreSQL)
type IBenchmarkRepository interface {
	Create(ctx context.Context, schema *modelBenchmark.BenchmarkSchemaPostgreSQL) error
	// List retorna todos os schemas (inclui schema_body)
	List(ctx context.Context) ([]*modelBenchmark.BenchmarkSchemaPostgreSQL, error)
	// GetByID busca um schema pelo ID (int64)
	GetByID(ctx context.Context, id int64) (*modelBenchmark.BenchmarkSchemaPostgreSQL, error)
	// Delete remove o registro
	Delete(ctx context.Context, id *int64) error
}

