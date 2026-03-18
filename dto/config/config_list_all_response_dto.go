package config

// ConfigListAllResponseDTO represents a collection of variable configurations
type ConfigListAllResponseDTO struct {
	Data []ConfigListResponseDTO `json:"data"`
}
