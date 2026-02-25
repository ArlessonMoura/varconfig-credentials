// Package config contains DTOs (Data Transfer Objects) for variable configuration operations
package config

// ConfigCreateRequestDTO represents the data required to create a new variable configuration
type ConfigCreateRequestDTO struct {
	Payload map[string]any `json:"payload" binding:"required"`
}
