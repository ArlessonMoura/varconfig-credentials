package benchmark

// BenchmarkCreateRequestDTO represents the data required to create a new benchmark schema
type BenchmarkCreateRequestDTO struct {
	Name   string                           `json:"name"`
	Schema map[string]SchemaFieldDefinition `json:"schema"`
}
