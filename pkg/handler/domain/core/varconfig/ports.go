package varconfig

import (
	"context"
	"projeto-crud-credencials/dto"
)

/*
IVarConfigService:
- O Handler chama estes métodos passando os DTOs de Request.
- O Service retorna os DTOs de Response.
*/
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
