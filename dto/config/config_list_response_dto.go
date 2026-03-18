package config

// ConfigListResponseDTO represents a lightweight item in configuration listing (without payload)
type ConfigListResponseDTO struct {
	ID          string `json:"id"`           // Unique identifier (SK)
	OrgID       string `json:"org_id"`       // Organization identifier (PK)
	BenchmarkID string `json:"benchmark_id"` // Benchmark identifier (part of PK)
	CreatedAt   string `json:"created_at"`   // Creation timestamp
	UpdatedAt   string `json:"updated_at"`   // Last update timestamp
}
