package benchmark_schema

import (
	"context"
	"fmt"
	"strconv"
	"time"

	repository "projeto-crud-credencials/internal/service"
	"projeto-crud-credencials/pkg/models"
)

type Service struct {
	relationalRepo repository.IRelationalRepository
	nosqlRepo      repository.INoSQLRepository
}

func NewService(relRepo repository.IRelationalRepository, noSqlRepo repository.INoSQLRepository) *Service {
	return &Service{
		relationalRepo: relRepo,
		nosqlRepo:      noSqlRepo,
	}
}

func (s *Service) RegisterSchema(ctx context.Context, name string, schemaBody map[string]any) (*models.BenchmarkSchemaRelational, error) {
	// 1. Prepara o modelo Relacional
	relationalModel := &models.BenchmarkSchemaRelational{
		Name: name,
	}

	// 2. Persiste no SQL (o ID será preenchido pelo banco no ponteiro)
	if err := s.relationalRepo.Create(ctx, relationalModel); err != nil {
		return nil, err
	}

	// 3. Prepara o modelo NoSQL usando o ID gerado pelo SQL
	noSqlModel := &models.BenchmarkSchemaNoSQL{
		ID:         strconv.FormatInt(relationalModel.ID, 10),
		SchemaBody: schemaBody,
		CreatedAt:  relationalModel.CreatedAt.Format(time.RFC3339),
	}

	// 4. Salva no DynamoDB
	if err := s.nosqlRepo.Save(ctx, noSqlModel); err != nil {
		// Rollback Compensatório: Deleta do SQL se o Dynamo falhar
		_ = s.relationalRepo.Delete(ctx, relationalModel.ID)
		return nil, fmt.Errorf("nosql storage failed, rolled back relational: %w", err)
	}

	return relationalModel, nil
}