package config

// ConfigUpdateRequestDTO represents the data required to update an existing variable configuration
type ConfigUpdateRequestDTO struct {
	Payload map[string]any `json:"payload" binding:"required"`
}
