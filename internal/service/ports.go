package service

import (
	"context"
	benchmarkModels "projeto-crud-credentials/pkg/models/benchmarkschema"
	varConfigModels "projeto-crud-credentials/pkg/models/varconfig"
)

// IVarConfigRepository define o contrato de persistência.
// O Service chama estes métodos passando as chaves já formatadas (PK e SK).
type IVarConfigRepository interface {
	// CREATE (PutItem do DynamoDB)
	Create(ctx context.Context, item *varConfigModels.VarConfigItem) error

	// List retorna apenas campos essenciais para listagem (sem payload pesado)
	List(ctx context.Context, pk string) ([]*varConfigModels.VarConfigListItem, error)
	
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
	Create(ctx context.Context, schema *benchmarkModels.BenchmarkSchemaPostgreSQL) error
	// Delete remove o registro em caso de falha no Dynamo (Rollback)
	Delete(ctx context.Context, id *int64) error
}

// INoSQLRepository define o contrato para persistência NoSQL genérica
type INoSQLRepository interface {
	Create(ctx context.Context, item *benchmarkModels.BenchmarkSchemaDynamoDB) error
	
	// List retorna apenas campos essenciais para listagem (sem schema_body pesado)
	List(ctx context.Context) ([]*benchmarkModels.BenchmarkSchemaListDynamoDB, error)
			
	//Na verdade lista pelo PK
	GetByID(ctx context.Context, id string) (*benchmarkModels.BenchmarkSchemaDynamoDB, error)
	Update(ctx context.Context, item *benchmarkModels.BenchmarkSchemaDynamoDB) error
	Delete(ctx context.Context, id string) error
}
