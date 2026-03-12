package benchmark

import (
	"errors"
	"strings"
)

// BenchmarkCreateRequestDTO represents the data required to create a new benchmark schema
type BenchmarkCreateRequestDTO struct {
	Name    string                           `json:"name"`
	Version string                           `json:"version"`
	Schema  map[string]SchemaFieldDefinition `json:"schema"`
}

func (r *BenchmarkCreateRequestDTO) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(r.Version) == "" {
		return errors.New("version is required")
	}
	return nil
}
