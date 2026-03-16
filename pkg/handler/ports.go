package handler

import (
	"context"
	dtoBenchmark "projeto-crud-credentials/dto/benchmark"
	dtoConfig "projeto-crud-credentials/dto/config"
	"projeto-crud-credentials/pkg/models"
)

type IVarConfigService interface {
	//CREATE
	Create(ctx context.Context, orgID string, benchmarkID string, input *dtoConfig.ConfigCreateRequestDTO) (*dtoConfig.ConfigResponseDTO, error)

	//READ
	List(ctx context.Context, orgID string, benchmarkID string) (*dtoConfig.ConfigListAllResponseDTO, error)
	GetByID(ctx context.Context, orgID string, benchmarkID string, id string) (*dtoConfig.ConfigResponseDTO, error)

	//UPDATE
	Update(ctx context.Context, orgID string, benchmarkID string, id string, input *dtoConfig.ConfigUpdateRequestDTO) (*dtoConfig.ConfigResponseDTO, error)

	//DELETE
	Delete(ctx context.Context, orgID string, benchmarkID string, id string) error
}

type IBenchmarkSchemaService interface {
	//CREATE
	Create(ctx context.Context, schemaRequest *dtoBenchmark.BenchmarkCreateRequestDTO) (*models.BenchmarkSchema, error)

	//READ
	List(ctx context.Context) (*dtoBenchmark.BenchmarkListResponseDTO, error)
	GetByID(ctx context.Context, id string) (*dtoBenchmark.BenchmarkResponseDTO, error)

	//DELETE
	Delete(ctx context.Context, id string) error
}
