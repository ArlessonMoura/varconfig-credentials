// Package benchmark contains DTOs (Data Transfer Objects) for benchmark schema operations
package benchmark

// SchemaFieldType defines the valid types for schema fields
type SchemaFieldType string

const (
	FieldTypeString  SchemaFieldType = "string"
	FieldTypeNumber  SchemaFieldType = "number"
	FieldTypeBoolean SchemaFieldType = "boolean"
	FieldTypeArray   SchemaFieldType = "array"
)

// ArrayFieldDefinition defines the structure for array-type fields
type ArrayFieldDefinition struct {
	Type SchemaFieldType `json:"type"`
}

// SchemaFieldDefinition represents a single field definition in a benchmark schema
type SchemaFieldDefinition struct {
	Type        SchemaFieldType       `json:"type"`
	Items       *ArrayFieldDefinition `json:"items,omitempty"` // Required when Type == "array"
	Required    bool                  `json:"required,omitempty"`
	Description string                `json:"description,omitempty"`
}
