package handler

import (
	"context"
	dtoBenchmark "projeto-crud-credentials/dto/benchmark"
	dtoConfig "projeto-crud-credentials/dto/config"
)


type IVarConfigService interface {
	//CREATE	
	Create(ctx context.Context, orgID string, benchmarkID string, input *dtoConfig.CreateVarConfigRequest) (*dtoConfig.VarConfigResponse, error)

	//READ
	List(ctx context.Context, orgID string, benchmarkID string) (*dtoConfig.ListVarConfigsResponse, error)
	GetByID(ctx context.Context, orgID string, benchmarkID string, id string) (*dtoConfig.VarConfigResponse, error)

	//UPDATE
	Update(ctx context.Context, orgID string, benchmarkID string, id string, input *dtoConfig.UpdateVarConfigRequest) (*dtoConfig.VarConfigResponse, error)

	//DELETE
	Delete(ctx context.Context, orgID string, benchmarkID string, id string) error
}

type IBenchmarkSchemaService interface {
	//CREATE
	Create(ctx context.Context, name string, schemaRequest *dtoBenchmark.CreateBenchmarkSchemaRequest) (dtoBenchmark.CreateBenchmarkSchemaResponse, error)
	
	//READ
	List(ctx context.Context) (*dtoBenchmark.ListBenchmarkSchemasResponse, error)
	GetByID(ctx context.Context, id string) (*dtoBenchmark.BenchmarkSchemaResponse, error)
	}
