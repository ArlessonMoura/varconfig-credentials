package helpers

import (
	"encoding/json"
	"fmt"
	"time"

	"projeto-crud-credentials/dto/benchmark"
	"projeto-crud-credentials/dto/config"
	"projeto-crud-credentials/pkg/models"
)

// MapVarConfigToResponse converte VarConfig para ConfigResponseDTO
func MapVarConfigToResponse(item *models.VarConfig) *config.ConfigResponseDTO {
	var payload map[string]any
	if item.Payload != nil {
		if err := json.Unmarshal(item.Payload, &payload); err != nil {
			payload = nil
		}
	}

	return &config.ConfigResponseDTO{
		ID:          fmt.Sprintf("%d", item.ID),
		Name:        item.Name,
		OrgID:       item.OrgID,
		BenchmarkID: item.BenchmarkID,
		Payload:     payload,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	}
}

// MapBenchmarkSchemaToResponse converte BenchmarkSchema para BenchmarkResponseDTO (sem schema)
func MapBenchmarkSchemaToResponse(item *models.BenchmarkSchema) *benchmark.BenchmarkResponseDTO {
	return &benchmark.BenchmarkResponseDTO{
		ID:         fmt.Sprintf("%d", item.ID),
		Name:       item.Name,
		Version:    item.Version,
		SchemaBody: nil,
		CreatedAt:  item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  item.UpdatedAt.Format(time.RFC3339),
	}
}

// MapBenchmarkSchemaToResponseWithSchema converte BenchmarkSchema para BenchmarkResponseDTO (com schema)
func MapBenchmarkSchemaToResponseWithSchema(item *models.BenchmarkSchema) *benchmark.BenchmarkResponseDTO {
	schemaBody := ConvertSchemaBodyToStringMap(item.Schema)

	return &benchmark.BenchmarkResponseDTO{
		ID:         fmt.Sprintf("%d", item.ID),
		Name:       item.Name,
		Version:    item.Version,
		SchemaBody: schemaBody,
		CreatedAt:  item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  item.UpdatedAt.Format(time.RFC3339),
	}
}

// ConvertSchemaBodyToStringMap converte json.RawMessage para map[string]string
func ConvertSchemaBodyToStringMap(schema json.RawMessage) map[string]string {
	var schemaDef map[string]benchmark.SchemaFieldDefinition
	if err := json.Unmarshal(schema, &schemaDef); err != nil {
		return nil
	}

	result := make(map[string]string)
	for k, v := range schemaDef {
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

// ConvertSchemaMapToStringMap converte map[string]SchemaFieldDefinition para map[string]string
func ConvertSchemaMapToStringMap(schema map[string]benchmark.SchemaFieldDefinition) map[string]string {
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
