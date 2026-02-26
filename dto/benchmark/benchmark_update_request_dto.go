package benchmark

// BenchmarkUpdateRequestDTO represents the data required to update an existing benchmark schema
type BenchmarkUpdateRequestDTO struct {
	Name   string                           `json:"name"`
	Schema map[string]SchemaFieldDefinition `json:"schema"`
}
