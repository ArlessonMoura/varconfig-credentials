// Package service provides domain services for the application.
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



// -------------------------------------------------


// IRelationalRepository define o contrato para persistência relacional (PostgreSQL)
type IRelationalRepository interface {
	Create(ctx context.Context, schema *models.BenchmarkSchemaRelational) error
	// Delete remove o registro em caso de falha no Dynamo (Rollback)
	Delete(ctx context.Context, id int64) error
}

// INoSQLRepository define o contrato para persistência NoSQL genérica
type INoSQLRepository interface {
	Save(ctx context.Context, item *models.BenchmarkSchemaNoSQL) error
	Get(ctx context.Context, id string) (*models.BenchmarkSchemaNoSQL, error)
	List(ctx context.Context) ([]models.BenchmarkSchemaNoSQL, error)
	Delete(ctx context.Context, id string) error
}
