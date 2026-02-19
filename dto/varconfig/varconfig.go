package varconfig

// CreateVarConfigRequest represents the data required to create a new variable configuration
type CreateVarConfigRequest struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

// UpdateVarConfigRequest represents the data required to update an existing variable configuration
type UpdateVarConfigRequest struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

// VarConfigResponse represents the complete variable configuration data returned from API
type VarConfigResponse struct {
	ID          string         `json:"id"`           // Unique identifier (SK)
	OrgID       string         `json:"org_id"`       // Organization identifier (PK)
	BenchmarkID string         `json:"benchmark_id"` // Benchmark identifier (part of PK)
	Payload     map[string]any `json:"payload"`     // Configuration data
	CreatedAt   string         `json:"created_at"`   // Creation timestamp
	UpdatedAt   string         `json:"updated_at"`   // Last update timestamp
}

// ListVarConfigsResponse represents a collection of variable configurations
type ListVarConfigsResponse struct {
	Data []VarConfigResponse `json:"data"`
}
