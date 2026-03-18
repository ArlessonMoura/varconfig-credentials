package config

import (
	"projeto-crud-credentials/internal/common/validation"
)

// ConfigCreateRequestDTO represents the data required to create a new variable configuration
type ConfigCreateRequestDTO struct {
	Name        string         `json:"name"`
	Payload     map[string]any `json:"payload" binding:"required"`
}

func (r *ConfigCreateRequestDTO) Validate() error {
	return validation.ValidateName(r.Name)
}
