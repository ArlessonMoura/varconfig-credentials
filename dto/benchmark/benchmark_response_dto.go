package benchmark

// BenchmarkResponseDTO represents the complete benchmark schema data returned from API
type BenchmarkResponseDTO struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	SchemaBody map[string]string `json:"schema_body"`
	CreatedAt  string            `json:"created_at"`
}
