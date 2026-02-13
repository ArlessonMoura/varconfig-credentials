package benchmark_schema_dto

type ValidSchemaType string

const (
	TypeString   ValidSchemaType = "string"
	TypeNumber   ValidSchemaType = "number"
	TypeBoolean  ValidSchemaType = "boolean"
	TypeArray    ValidSchemaType = "array"
)

type ArrayItemsType struct {
	Type ValidSchemaType `json:"type"`
}

type SchemaProperty struct {
	Type        ValidSchemaType `json:"type"`
	Items       *ArrayItemsType `json:"items,omitempty"` // Obrigatório se Type == "array"
	Required    bool            `json:"required,omitempty"`
	Description string          `json:"description,omitempty"`
}

type InternalRegisterSchemaRequest struct {
	Name   string                   `json:"name"`
	Schema map[string]SchemaProperty `json:"schema"`
}

// DTOs de resposta para o Handler
type BenchmarkSchemaResponse struct {
	ID         string                   `json:"id"`
	Name       string                   `json:"name"`
	SchemaBody map[string]string        `json:"schema_body"`
	CreatedAt  string                   `json:"created_at"`
}

type RegisterSchemaResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type ListBenchmarkSchemasResponse struct {
	Data  []BenchmarkSchemaResponse `json:"data"`
	Count int                      `json:"count"`
}
