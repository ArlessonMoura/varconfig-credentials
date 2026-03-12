// Package config contains DTOs (Data Transfer Objects) for variable configuration operations
package config

import (
	"errors"
	"strings"
)

// ConfigCreateRequestDTO represents the data required to create a new variable configuration
type ConfigCreateRequestDTO struct {
	Name        string         `json:"name"`
	Payload     map[string]any `json:"payload" binding:"required"`
}

func (r *ConfigCreateRequestDTO) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}
