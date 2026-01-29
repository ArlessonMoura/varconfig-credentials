package benchmark_schema_dto

type InternalRegisterSchemaRequest struct {
	Name   string         `json:"name"`
	Schema map[string]string `json:"schema"`
}