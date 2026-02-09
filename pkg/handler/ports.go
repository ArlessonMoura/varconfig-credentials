// Package handler provides HTTP handlers for the application.
package handler

import (
	"context"
	dtoSchema "projeto-crud-credencials/dto/benchmark_schema_dto"
	dto "projeto-crud-credencials/dto/varconfig_dto"
)


type IVarConfigService interface {
	// List retorna o wrapper de data contendo a lista
	List(ctx context.Context, orgID string, benchmarkID string) (dto.ListVarConfigResponse, error)

	// GetByID retorna uma única resposta formatada
	GetByID(ctx context.Context, orgID string, benchmarkID string, id string) (dto.VarConfigResponse, error)

	// Create recebe o payload do request e retorna o ID gerado ou o objeto completo
	Create(ctx context.Context, orgID string, benchmarkID string, input dto.CreateVarConfigRequest) (dto.VarConfigResponse, error)

	// Update recebe o novo payload e as chaves de busca (PK/SK)
	Update(ctx context.Context, orgID string, benchmarkID string, id string, input dto.UpdateVarConfigRequest) (dto.VarConfigResponse, error)

	// Delete remove o registro
	Delete(ctx context.Context, orgID string, benchmarkID string, id string) error
}

type IBenchmarkSchemaService interface {
	// RegisterSchema registra um novo schema nos bancos relacional e NoSQL
	RegisterSchema(ctx context.Context, name string, schemaRequest *dtoSchema.InternalRegisterSchemaRequest) (dtoSchema.RegisterSchemaResponse, error)
	
	// GetSchemaByID recupera um schema específico pelo ID
	GetSchemaByID(ctx context.Context, id string) (dtoSchema.BenchmarkSchemaResponse, error)
	
	// ListAllSchemas recupera todos os schemas disponíveis
	ListAllSchemas(ctx context.Context) (dtoSchema.ListBenchmarkSchemasResponse, error)
}
