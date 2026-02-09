package benchmark_schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	dto "projeto-crud-credencials/dto/benchmark_schema_dto"
	repository "projeto-crud-credencials/internal/service"
	"projeto-crud-credencials/pkg/models"
)

var (
ErrInvalidInput = errors.New("invalid input: name is required")
ErrEmptySchema  = errors.New("schema cannot be empty")
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

// RegisterSchema registra um novo schema de benchmark de forma atômica
// 1. Insere o registro no banco relacional (PostgreSQL) e obtém o ID
// 2. Armazena o corpo completo no DynamoDB usando o ID como chave
// 3. Se o DynamoDB falhar, executa rollback do banco relacional (transação compensatória)

func (s *Service) RegisterSchema(
ctx context.Context,
name string,
schemaRequest *dto.InternalRegisterSchemaRequest,
) (*models.BenchmarkSchemaRelational, error) {
// Validações
if name == "" {
return nil, ErrInvalidInput
}

if len(schemaRequest.Schema) == 0 {
return nil, ErrEmptySchema
}

// 1. Registrar no banco relacional (PostgreSQL)
relationalModel := &models.BenchmarkSchemaRelational{
Name: name,
}

if err := s.relationalRepo.Create(ctx, relationalModel); err != nil {
return nil, fmt.Errorf("failed to create benchmark schema in relational database: %w", err)
}

noSqlModel := &models.BenchmarkSchemaNoSQL{
ID:         strconv.FormatInt(relationalModel.ID, 10),
SchemaBody: convertSchemaBodyToStringMap(schemaRequest.Schema),
CreatedAt:  relationalModel.CreatedAt.Format(time.RFC3339),
}

if err := s.nosqlRepo.Save(ctx, noSqlModel); err != nil {
// 3. Rollback Compensatório: Deleta do SQL se o DynamoDB falhar
deleteErr := s.relationalRepo.Delete(ctx, relationalModel.ID)
if deleteErr != nil {
fmt.Printf("CRITICAL: Failed to rollback schema %d after DynamoDB failure: %v\n", relationalModel.ID, deleteErr)
return nil, fmt.Errorf("failed to save schema in DynamoDB and failed to rollback: %w, rollback error: %w", err, deleteErr)
}
return nil, fmt.Errorf("failed to save schema in NoSQL database, rolled back relational entry: %w", err)
}

return relationalModel, nil
}

func convertSchemaBodyToStringMap(schema map[string]dto.SchemaProperty) map[string]string {
	result := make(map[string]string)
	for k, v := range schema {
		// transforma a struct SchemaProperty em uma string JSON válida
		jsonData, err := json.Marshal(v)
		if err != nil {
			result[k] = fmt.Sprintf(`{"type":"%s","error":"marshal_failed"}`, v.Type)
			continue
		}
		result[k] = string(jsonData)
	}
	return result
}
