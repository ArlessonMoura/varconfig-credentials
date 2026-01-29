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
