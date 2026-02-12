package handler

import (
	"context"
	dtoSchema "projeto-crud-credencials/dto/benchmark_schema_dto"
	dto "projeto-crud-credencials/dto/varconfig_dto"
)


type IVarConfigService interface {
	List(ctx context.Context, orgID string, benchmarkID string) (dto.ListVarConfigResponse, error)

	GetByID(ctx context.Context, orgID string, benchmarkID string, id string) (dto.VarConfigResponse, error)

	// Create recebe o payload do request e retorna o ID gerado ou o objeto completo
	Create(ctx context.Context, orgID string, benchmarkID string, input dto.CreateVarConfigRequest) (dto.VarConfigResponse, error)

	Update(ctx context.Context, orgID string, benchmarkID string, id string, input dto.UpdateVarConfigRequest) (dto.VarConfigResponse, error)

	Delete(ctx context.Context, orgID string, benchmarkID string, id string) error
}

type IBenchmarkSchemaService interface {
	Create(ctx context.Context, name string, schemaRequest *dtoSchema.InternalRegisterSchemaRequest) (dtoSchema.RegisterSchemaResponse, error)
	
	GetByID(ctx context.Context, id string) (dtoSchema.BenchmarkSchemaResponse, error)
	
	List(ctx context.Context) (dtoSchema.ListBenchmarkSchemasResponse, error)
}
