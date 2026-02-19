package handler

import (
	"context"
	dtoSchema "projeto-crud-credencials/dto/benchmarkschema"
	dto "projeto-crud-credencials/dto/varconfig"
)


type IVarConfigService interface {
	//CREATE	
	Create(ctx context.Context, orgID string, benchmarkID string, input *dto.VarConfigCreationPayload) (*dto.VarConfigData, error)

	//READ
	List(ctx context.Context, orgID string, benchmarkID string) (*dto.VarConfigCollectionResponse, error)
	GetByID(ctx context.Context, orgID string, benchmarkID string, id string) (*dto.VarConfigData, error)

	//UPDATE
	Update(ctx context.Context, orgID string, benchmarkID string, id string, input *dto.VarConfigUpdatePayload) (*dto.VarConfigData, error)

	//DELETE
	Delete(ctx context.Context, orgID string, benchmarkID string, id string) error
}

type IBenchmarkSchemaService interface {
	//CREATE
	Create(ctx context.Context, name string, schemaRequest *dtoSchema.BenchmarkSchemaCreationRequest) (dtoSchema.BenchmarkSchemaRegistrationResponse, error)
	
	//READ
	GetByID(ctx context.Context, id string) (*dtoSchema.BenchmarkSchemaDetails, error)	
	List(ctx context.Context) (*dtoSchema.BenchmarkSchemaCollectionResponse, error)
	
	//====================
	// 
	//====================
	
	//UPDATE
	Update(ctx context.Context, id string, name string, schemaRequest *dtoSchema.BenchmarkSchemaCreationRequest) (*dtoSchema.BenchmarkSchemaDetails, error)
	
	//DELETE
	Delete(ctx context.Context, id string) error
}
