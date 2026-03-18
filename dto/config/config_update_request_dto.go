package config

import (
	"projeto-crud-credentials/internal/common/validation"
)

// ConfigUpdateRequestDTO represents the data required to update an existing variable configuration
type ConfigUpdateRequestDTO struct {
	Name    string         `json:"name"`
	Payload map[string]any `json:"payload" binding:"required"`
}

func (r *ConfigUpdateRequestDTO) Validate() error {
	return validation.ValidateName(r.Name)
}
