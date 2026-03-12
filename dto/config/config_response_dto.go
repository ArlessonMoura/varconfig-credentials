package config

// ConfigResponseDTO represents the complete variable configuration data returned from API
type ConfigResponseDTO struct {
	ID          string         `json:"id"`           // Unique identifier (SK)
	Name        string         `json:"name"`         // Configuration name
	OrgID       string         `json:"org_id"`       // Organization identifier (PK)
	BenchmarkID string         `json:"benchmark_id"` // Benchmark identifier (part of PK)
	Payload     map[string]any `json:"payload"`      // Configuration data
	CreatedAt   string         `json:"created_at"`   // Creation timestamp
	UpdatedAt   string         `json:"updated_at"`   // Last update timestamp
}
