package service

import (
	"context"
	"projeto-crud-credencials/pkg/models"
)

// IVarConfigRepository define o contrato de persistência.
// O Service chama estes métodos passando as chaves já formatadas (PK e SK).
type IVarConfigRepository interface {
	// Save é usado tanto para Create quanto para Update (PutItem do DynamoDB)
	Save(ctx context.Context, item models.VarConfigItem) error

	// Get busca um item específico pela chave primária completa
	Get(ctx context.Context, pk string, sk string) (*models.VarConfigItem, error)

	// ListByPK busca todos os itens que compartilham a mesma partição (Benchmark)
	ListByPK(ctx context.Context, pk string) ([]models.VarConfigItem, error)

	// Delete remove o item via chaves compostas
	Delete(ctx context.Context, pk string, sk string) error
}
