package benchmark

// BenchmarkListResponseDTO represents a collection of benchmark schemas
type BenchmarkListResponseDTO struct {
	Data  []BenchmarkResponseDTO `json:"data"`
	Count int                    `json:"count"`
}
