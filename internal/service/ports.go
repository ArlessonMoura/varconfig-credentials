package service

import (
	"context"
	benchmarkModels "projeto-crud-credencials/pkg/models/benchmarkschema"
	varConfigModels "projeto-crud-credencials/pkg/models/varconfig"
)

// IVarConfigRepository define o contrato de persistência.
// O Service chama estes métodos passando as chaves já formatadas (PK e SK).
type IVarConfigRepository interface {
	// CREATE (PutItem do DynamoDB)
	Create(ctx context.Context, item *varConfigModels.VarConfigItem) error

	//TODO: observar se a List esta sendo usada corretamente para listar todos os items
	List(ctx context.Context, pk string) ([]*varConfigModels.VarConfigItem, error)
	// GetByID busca um item específico pela PK e SK
	GetByID(ctx context.Context, pk string, sk string) (*varConfigModels.VarConfigItem, error)

	// UPDATE (PutItem do DynamoDB)
	Update(ctx context.Context, item *varConfigModels.VarConfigItem) error

	// Delete remove o item via chaves compostas
	Delete(ctx context.Context, pk string, sk string) error
}



// -------------------------------------------------


// IRelationalRepository define o contrato para persistência relacional (PostgreSQL)
type IRelationalRepository interface {
	Create(ctx context.Context, schema *benchmarkModels.BenchmarkSchemaRelational) error
	// Delete remove o registro em caso de falha no Dynamo (Rollback)
	Delete(ctx context.Context, id *int64) error
}

// INoSQLRepository define o contrato para persistência NoSQL genérica
type INoSQLRepository interface {
	Create(ctx context.Context, item *benchmarkModels.BenchmarkSchemaNoSQL) error
	
	//Na verdade lista pelo PK
	GetByID(ctx context.Context, id string) (*benchmarkModels.BenchmarkSchemaNoSQL, error)
	List(ctx context.Context) ([]*benchmarkModels.BenchmarkSchemaNoSQL, error)
	Update(ctx context.Context, item *benchmarkModels.BenchmarkSchemaNoSQL) error
	Delete(ctx context.Context, id string) error
}
