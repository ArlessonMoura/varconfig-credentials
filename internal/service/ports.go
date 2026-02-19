package service

import (
	"context"
	modelBenchmark "projeto-crud-credentials/pkg/models/benchmark"
	modelConfig "projeto-crud-credentials/pkg/models/config"
)

// IVarConfigRepository define o contrato de persistência.
// O Service chama estes métodos passando as chaves já formatadas (PK e SK).
type IVarConfigRepository interface {
	// CREATE (PutItem do DynamoDB)
	Create(ctx context.Context, item *modelConfig.VarConfigItem) error

	// List retorna apenas campos essenciais para listagem (sem payload pesado)
	List(ctx context.Context, pk string) ([]*modelConfig.VarConfigListItem, error)
	
	// GetByID busca um item específico pela PK e SK
	GetByID(ctx context.Context, pk string, sk string) (*modelConfig.VarConfigItem, error)

	// UPDATE (PutItem do DynamoDB)
	Update(ctx context.Context, item *modelConfig.VarConfigItem) error

	// Delete remove o item via chaves compostas
	Delete(ctx context.Context, pk string, sk string) error
}



// -------------------------------------------------


// IRelationalRepository define o contrato para persistência relacional (PostgreSQL)
type IRelationalRepository interface {
	Create(ctx context.Context, schema *modelBenchmark.BenchmarkSchemaPostgreSQL) error
	// Delete remove o registro em caso de falha no Dynamo (Rollback)
	Delete(ctx context.Context, id *int64) error
}

// INoSQLRepository define o contrato para persistência NoSQL genérica
type INoSQLRepository interface {
	Create(ctx context.Context, item *modelBenchmark.BenchmarkSchemaDynamoDB) error
	
	// List retorna apenas campos essenciais para listagem (sem schema_body pesado)
	List(ctx context.Context) ([]*modelBenchmark.BenchmarkSchemaListDynamoDB, error)			
	//Na verdade lista pelo PK
	GetByID(ctx context.Context, id string) (*modelBenchmark.BenchmarkSchemaDynamoDB, error)


// 	// ===========

// 	// ===========
// 	Update(ctx context.Context, item *modelBenchmark.BenchmarkSchemaDynamoDB) error
// 	Delete(ctx context.Context, id string) error
}
