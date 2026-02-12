package handler

import (
	"context"
	dtoSchema "projeto-crud-credencials/dto/benchmark_schema"
	dto "projeto-crud-credencials/dto/varconfig"
)


type IVarConfigService interface {
	//CREATE	
	Create(ctx context.Context, orgID string, benchmarkID string, input dto.CreateVarConfigRequest) (dto.VarConfigResponse, error)

	//READ
	List(ctx context.Context, orgID string, benchmarkID string) (dto.ListVarConfigResponse, error)
	GetByID(ctx context.Context, orgID string, benchmarkID string, id string) (dto.VarConfigResponse, error)

	//UPDATE
	Update(ctx context.Context, orgID string, benchmarkID string, id string, input dto.UpdateVarConfigRequest) (dto.VarConfigResponse, error)

	//DELETE
	Delete(ctx context.Context, orgID string, benchmarkID string, id string) error
}

type IBenchmarkSchemaService interface {
	//CREATE
	Create(ctx context.Context, name string, schemaRequest *dtoSchema.InternalRegisterSchemaRequest) (dtoSchema.RegisterSchemaResponse, error)
	
	//READ
	GetByID(ctx context.Context, id string) (dtoSchema.BenchmarkSchemaResponse, error)	
	List(ctx context.Context) (dtoSchema.ListBenchmarkSchemasResponse, error)
	// Abaixo, ambos são implementações extra --
	//UPDATE
	Update(ctx context.Context, id string, name string, schemaRequest *dtoSchema.InternalRegisterSchemaRequest) (dtoSchema.BenchmarkSchemaResponse, error)
	
	//DELETE
	Delete(ctx context.Context, id string) error
}
