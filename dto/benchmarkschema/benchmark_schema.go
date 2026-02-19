package benchmark_schema

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

// BenchmarkSchemaCreationRequest represents the data required to create a new benchmark schema
type BenchmarkSchemaCreationRequest struct {
	Name   string                          `json:"name"`
	Schema map[string]SchemaFieldDefinition `json:"schema"`
}

// BenchmarkSchemaUpdateRequest represents the data required to update an existing benchmark schema
type BenchmarkSchemaUpdateRequest struct {
	Name   string                          `json:"name"`
	Schema map[string]SchemaFieldDefinition `json:"schema"`
}

// BenchmarkSchemaDetails represents the complete benchmark schema data returned from API
type BenchmarkSchemaDetails struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	SchemaBody map[string]string `json:"schema_body"`
	CreatedAt  string            `json:"created_at"`
}

// BenchmarkSchemaRegistrationResponse represents the response after successful schema registration
type BenchmarkSchemaRegistrationResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// BenchmarkSchemaCollectionResponse represents a collection of benchmark schemas
type BenchmarkSchemaCollectionResponse struct {
	Data  []BenchmarkSchemaDetails `json:"data"`
	Count int                     `json:"count"`
}
