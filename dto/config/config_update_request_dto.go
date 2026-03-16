package config

import (
	"errors"
	"strings"
)

// ConfigUpdateRequestDTO represents the data required to update an existing variable configuration
type ConfigUpdateRequestDTO struct {
	Name    string         `json:"name"`
	Payload map[string]any `json:"payload" binding:"required"`
}

func (r *ConfigUpdateRequestDTO) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}
