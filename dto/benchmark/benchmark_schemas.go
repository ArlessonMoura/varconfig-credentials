package benchmark

// SchemaFieldType defines the valid types for schema fields
type SchemaFieldType string

const (
	FieldTypeString   SchemaFieldType = "string"
	FieldTypeNumber   SchemaFieldType = "number"
	FieldTypeBoolean  SchemaFieldType = "boolean"
	FieldTypeArray    SchemaFieldType = "array"
)

// ArrayFieldDefinition defines the structure for array-type fields
type ArrayFieldDefinition struct {
	Type SchemaFieldType `json:"type"`
}

// SchemaFieldDefinition represents a single field definition in a benchmark schema
type SchemaFieldDefinition struct {
	Type        SchemaFieldType      `json:"type"`
	Items       *ArrayFieldDefinition `json:"items,omitempty"` // Required when Type == "array"
	Required    bool                `json:"required,omitempty"`
	Description string              `json:"description,omitempty"`
}

// CreateBenchmarkSchemaRequest represents the data required to create a new benchmark schema
type CreateBenchmarkSchemaRequest struct {
	Name   string                          `json:"name"`
	Schema map[string]SchemaFieldDefinition `json:"schema"`
}

// UpdateBenchmarkSchemaRequest represents the data required to update an existing benchmark schema
type UpdateBenchmarkSchemaRequest struct {
	Name   string                          `json:"name"`
	Schema map[string]SchemaFieldDefinition `json:"schema"`
}

// BenchmarkSchemaResponse represents the complete benchmark schema data returned from API
type BenchmarkSchemaResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	SchemaBody map[string]string `json:"schema_body"`
	CreatedAt  string            `json:"created_at"`
}

// CreateBenchmarkSchemaResponse represents the response after successful schema registration
type CreateBenchmarkSchemaResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// ListBenchmarkSchemasResponse represents a collection of benchmark schemas
type ListBenchmarkSchemasResponse struct {
	Data  []BenchmarkSchemaResponse `json:"data"`
	Count int                     `json:"count"`
}
