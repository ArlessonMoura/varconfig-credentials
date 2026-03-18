package benchmark

import (
	"projeto-crud-credentials/internal/common/validation"
)

// BenchmarkCreateRequestDTO represents the data required to create a new benchmark schema
type BenchmarkCreateRequestDTO struct {
	Name    string                           `json:"name"`
	Version string                           `json:"version"`
	Schema  map[string]SchemaFieldDefinition `json:"schema"`
}

func (r *BenchmarkCreateRequestDTO) Validate() error {
	if err := validation.ValidateName(r.Name); err != nil {
		return err
	}
	return validation.ValidateRequiredField(r.Version, "version")
}
