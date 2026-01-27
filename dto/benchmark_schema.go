package dto

type InternalRegisterSchemaRequest struct {
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema"`
}