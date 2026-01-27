package service

import (
	"context"
	"projeto-crud-credencials/pkg/models"
)

// IVarConfigRepository define o contrato de persistência.
// O Service chama estes métodos passando as chaves já formatadas (PK e SK).
type IVarConfigRepository interface {
	// Save é usado tanto para Create quanto para Update (PutItem do DynamoDB)
	Save(ctx context.Context, item models.VarConfigItem) error

	// Get busca um item específico pela chave primária completa
	Get(ctx context.Context, pk string, sk string) (*models.VarConfigItem, error)

	// ListByPK busca todos os itens que compartilham a mesma partição (Benchmark)
	ListByPK(ctx context.Context, pk string) ([]models.VarConfigItem, error)

	// Delete remove o item via chaves compostas
	Delete(ctx context.Context, pk string, sk string) error
}


// ======================================


// IBenchmarkSchemaService define a lógica interna de coordenação
type IBenchmarkSchemaService interface {
	RegisterSchema(ctx context.Context, name string, schemaBody map[string]any) (*models.BenchmarkSchemaRelational, error)
}

// IRelationalRepository lida com o banco SQL (Postgres/MySQL)
type IRelationalRepository interface {
	// Create agora recebe o modelo e retorna ele preenchido (com ID e Timestamps)
	Create(ctx context.Context, schema *models.BenchmarkSchemaRelational) error
	// Delete remove o registro em caso de falha no Dynamo (Rollback)
	Delete(ctx context.Context, id int64) error
}

// INoSQLRepository lida com o DynamoDB
type INoSQLRepository interface {
	// Save agora recebe o modelo NoSQL completo
	Save(ctx context.Context, item *models.BenchmarkSchemaNoSQL) error
}
