package benchmark

// BenchmarkCreateResponseDTO represents the response after successful schema registration
type BenchmarkCreateResponseDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
