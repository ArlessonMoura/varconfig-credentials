package varconfig

// VarConfigCreationPayload represents the data required to create a new variable configuration
type VarConfigCreationPayload struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

// VarConfigUpdatePayload represents the data required to update an existing variable configuration
type VarConfigUpdatePayload struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

// VarConfigData represents the complete variable configuration data returned from API
type VarConfigData struct {
	ID          string         `json:"id"`           // Unique identifier (SK)
	OrgID       string         `json:"org_id"`       // Organization identifier (PK)
	BenchmarkID string         `json:"benchmark_id"` // Benchmark identifier (part of PK)
	Payload     map[string]any `json:"payload"`     // Configuration data
	CreatedAt   string         `json:"created_at"`   // Creation timestamp
	UpdatedAt   string         `json:"updated_at"`   // Last update timestamp
}

// VarConfigCollectionResponse represents a collection of variable configurations
type VarConfigCollectionResponse struct {
	Data []VarConfigData `json:"data"`
}
