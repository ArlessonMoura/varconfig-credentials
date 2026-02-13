// Package benchmark_schema provides services for managing benchmark schemas in both relational and NoSQL databases.
package benchmark_schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	dto "projeto-crud-credencials/dto/benchmarkschema"
	repository "projeto-crud-credencials/internal/service"
	models "projeto-crud-credencials/pkg/models/benchmark_schema"
)

var (
	ErrInvalidInput   = errors.New("invalid input: name is required")
	ErrEmptySchema    = errors.New("schema cannot be empty")
	ErrSchemaNotFound = errors.New("schema not found")
)

type Service struct {
	relationalRepo repository.IRelationalRepository
	nosqlRepo      repository.INoSQLRepository

}

func NewService(relRepo repository.IRelationalRepository, noSQLRepo repository.INoSQLRepository) *Service {
	return &Service{
		relationalRepo: relRepo,
		nosqlRepo:      noSQLRepo,
	}
}

func (s *Service) Create(ctx context.Context, name string, schemaRequest *dto.BenchmarkSchemaCreationRequest) (dto.BenchmarkSchemaRegistrationResponse, error) {
	// Validações
	if name == "" {
		return dto.BenchmarkSchemaRegistrationResponse{}, ErrInvalidInput
	}

	if len(schemaRequest.Schema) == 0 {
		return dto.BenchmarkSchemaRegistrationResponse{}, ErrEmptySchema
	}

	// 1. Registrar no banco relacional (PostgreSQL)
	relationalModel := &models.BenchmarkSchemaRelational{
		Name: name,
	}

	if err := s.relationalRepo.Create(ctx, relationalModel); err != nil {
		return dto.BenchmarkSchemaRegistrationResponse{}, fmt.Errorf("failed to create benchmark schema in relational database: %w", err)
	}

	noSQLModel := &models.BenchmarkSchemaNoSQL{
		ID:         strconv.FormatInt(relationalModel.ID, 10),
		SchemaBody: convertSchemaBodyToStringMap(schemaRequest.Schema),
		CreatedAt:  relationalModel.CreatedAt.Format(time.RFC3339),
	}

	if err := s.nosqlRepo.Create(ctx, noSQLModel); err != nil {
		// 3. Rollback Compensatório: Deleta do SQL se o DynamoDB falhar
		deleteErr := s.relationalRepo.Delete(ctx, relationalModel.ID)
		if deleteErr != nil {
			fmt.Printf("CRITICAL: Failed to rollback schema %d after DynamoDB failure: %v\n", relationalModel.ID, deleteErr)
			return dto.BenchmarkSchemaRegistrationResponse{}, fmt.Errorf("failed to save schema in DynamoDB and failed to rollback: %w, rollback error: %w", err, deleteErr)
		}
		return dto.BenchmarkSchemaRegistrationResponse{}, fmt.Errorf("failed to save schema in NoSQL database, rolled back relational entry: %w", err)
	}

	return dto.BenchmarkSchemaRegistrationResponse{
		ID:        noSQLModel.ID,
		Name:      relationalModel.Name,
		CreatedAt: relationalModel.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetByID recupera um schema específico pelo ID
func (s *Service) GetByID(ctx context.Context, id string) (dto.BenchmarkSchemaDetails, error) {
	if id == "" {
		return dto.BenchmarkSchemaDetails{}, ErrSchemaNotFound
	}

	item, err := s.nosqlRepo.GetByID(ctx, id)
	if err != nil {
		return dto.BenchmarkSchemaDetails{}, fmt.Errorf("failed to get schema from repository: %w", err)
	}

	if item == nil {
		return dto.BenchmarkSchemaDetails{}, ErrSchemaNotFound
	}

	return dto.BenchmarkSchemaDetails{
		ID:         item.ID,
		Name:       "", // Nome não está disponível no NoSQL, poderia buscar do relacional se necessário
		SchemaBody: item.SchemaBody,
		CreatedAt:  item.CreatedAt,
	}, nil
}

// List recupera todos os schemas disponíveis
func (s *Service) List(ctx context.Context) (dto.BenchmarkSchemaCollectionResponse, error) {
	items, err := s.nosqlRepo.List(ctx)
	if err != nil {
		return dto.BenchmarkSchemaCollectionResponse{}, fmt.Errorf("failed to list schemas from repository: %w", err)
	}

	if items == nil {
		return dto.BenchmarkSchemaCollectionResponse{Data: []dto.BenchmarkSchemaDetails{}, Count: 0}, nil
	}

	var data []dto.BenchmarkSchemaDetails
	for _, item := range items {
		data = append(data, dto.BenchmarkSchemaDetails{
			ID:         item.ID,
			Name:       "", // Nome não está disponível no NoSQL
			SchemaBody: item.SchemaBody,
			CreatedAt:  item.CreatedAt,
		})
	}

	return dto.BenchmarkSchemaCollectionResponse{
		Data:  data,
		Count: len(data),
	}, nil
}


// Update atualiza um schema existente usando dual-write
func (s *Service) Update(ctx context.Context, id string, name string, schemaRequest *dto.BenchmarkSchemaCreationRequest) (dto.BenchmarkSchemaDetails, error) {
	if id == "" || name == "" {
		return dto.BenchmarkSchemaDetails{}, ErrInvalidInput
	}

	if len(schemaRequest.Schema) == 0 {
		return dto.BenchmarkSchemaDetails{}, ErrEmptySchema
	}

	// 1. Buscar schema atual do NoSQL
	existing, err := s.nosqlRepo.GetByID(ctx, id)
	if err != nil {
		return dto.BenchmarkSchemaDetails{}, fmt.Errorf("failed to get existing schema: %w", err)
	}
	if existing == nil {
		return dto.BenchmarkSchemaDetails{}, ErrSchemaNotFound
	}

	// 2. Atualizar NoSQL com novo schema
	updatedNoSQL := &models.BenchmarkSchemaNoSQL{
		ID:         id,
		SchemaBody: convertSchemaBodyToStringMap(schemaRequest.Schema),
		CreatedAt:  existing.CreatedAt, // Mantém created_at original
	}

	if err := s.nosqlRepo.Update(ctx, updatedNoSQL); err != nil {
		return dto.BenchmarkSchemaDetails{}, fmt.Errorf("failed to update schema in NoSQL database: %w", err)
	}

	// Atualizar nome no banco relacional (se necessário)
	// Nota: Esta implementação assume que só o schema body muda, não o nome
	// Se o nome também precisar ser atualizado, precisaríamos de um método Update no relacional repo

	return dto.BenchmarkSchemaDetails{
		ID:         updatedNoSQL.ID,
		Name:       name, // Nome passado como parâmetro
		SchemaBody: updatedNoSQL.SchemaBody,
		CreatedAt:  updatedNoSQL.CreatedAt,
	}, nil
}

// Delete remove um schema usando dual-write
func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrSchemaNotFound
	}

	// 1. Verificar se schema existe no NoSQL
	existing, err := s.nosqlRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check schema existence: %w", err)
	}
	if existing == nil {
		return ErrSchemaNotFound
	}

	// 2. Converter ID para int64 para o relacional
	relationalID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid schema ID format: %w", err)
	}

	// 3. Remover do NoSQL
	if err := s.nosqlRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete schema from NoSQL database: %w", err)
	}

	// 4. Remover do relacional
	if err := s.relationalRepo.Delete(ctx, relationalID); err != nil {
		// Rollback compensatório: tentar recriar no NoSQL se falhar a deleção no relacional
		rollbackErr := s.nosqlRepo.Create(ctx, existing)
		if rollbackErr != nil {
			fmt.Printf("CRITICAL: Failed to rollback schema %s after relational delete failure: %v\n", id, rollbackErr)
			return fmt.Errorf("failed to delete from relational database and failed to rollback: %w, rollback error: %w", err, rollbackErr)
		}
		return fmt.Errorf("failed to delete from relational database, rolled back NoSQL entry: %w", err)
	}

	return nil
}

func convertSchemaBodyToStringMap(schema map[string]dto.SchemaFieldDefinition) map[string]string {
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
