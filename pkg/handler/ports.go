package handler

import (
	"context"
	dtoSchema "projeto-crud-credentials/dto/benchmarkschema"
	dto "projeto-crud-credentials/dto/varconfig"
)


type IVarConfigService interface {
	//CREATE	
	Create(ctx context.Context, orgID string, benchmarkID string, input *dto.CreateVarConfigRequest) (*dto.VarConfigResponse, error)

	//READ
	List(ctx context.Context, orgID string, benchmarkID string) (*dto.ListVarConfigsResponse, error)
	GetByID(ctx context.Context, orgID string, benchmarkID string, id string) (*dto.VarConfigResponse, error)

	//UPDATE
	Update(ctx context.Context, orgID string, benchmarkID string, id string, input *dto.UpdateVarConfigRequest) (*dto.VarConfigResponse, error)

	//DELETE
	Delete(ctx context.Context, orgID string, benchmarkID string, id string) error
}

type IBenchmarkSchemaService interface {
	//CREATE
	Create(ctx context.Context, name string, schemaRequest *dtoSchema.CreateBenchmarkSchemaRequest) (dtoSchema.CreateBenchmarkSchemaResponse, error)
	
	//READ
	GetByID(ctx context.Context, id string) (*dtoSchema.BenchmarkSchemaResponse, error)	
	List(ctx context.Context) (*dtoSchema.ListBenchmarkSchemasResponse, error)
	
	//====================
	// 
	//====================
	
	//UPDATE
	Update(ctx context.Context, id string, name string, schemaRequest *dtoSchema.UpdateBenchmarkSchemaRequest) (*dtoSchema.BenchmarkSchemaResponse, error)
	
	//DELETE
	Delete(ctx context.Context, id string) error
}
