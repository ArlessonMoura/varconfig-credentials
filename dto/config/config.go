package config

// VarConfigResponse represents the complete variable configuration data returned from API
type VarConfigResponse struct {
	ID          string         `json:"id"`           // Unique identifier (SK)
	OrgID       string         `json:"org_id"`       // Organization identifier (PK)
	BenchmarkID string         `json:"benchmark_id"` // Benchmark identifier (part of PK)
	Payload     map[string]any `json:"payload"`     // Configuration data
	CreatedAt   string         `json:"created_at"`   // Creation timestamp
	UpdatedAt   string         `json:"updated_at"`   // Last update timestamp
}

type CreateVarConfigRequest struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

type UpdateVarConfigRequest struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

type ListVarConfigsResponse struct {
	Data []VarConfigResponse `json:"data"`
}
