package benchmark

import (
	"fmt"
)

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

// Validate verifica se um valor está em conformidade com a definição de array
func (a *ArrayFieldDefinition) Validate(fieldName string, value any) error {
	switch a.Type {
	case FieldTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("array element in field '%s' must be string", fieldName)
		}
	case FieldTypeNumber:
		switch value.(type) {
		case float64, float32, int, int64, int32:
			// ok
		default:
			return fmt.Errorf("array element in field '%s' must be number", fieldName)
		}
	case FieldTypeBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("array element in field '%s' must be boolean", fieldName)
		}
	}
	return nil
}

// SchemaFieldDefinition represents a single field definition in a benchmark schema
type SchemaFieldDefinition struct {
	Type        SchemaFieldType       `json:"type"`
	Items       *ArrayFieldDefinition `json:"items,omitempty"` // Required when Type == "array"
	Required    bool                  `json:"required,omitempty"`
	Description string                `json:"description,omitempty"`
}

// Validate verifica se um valor está em conformidade com a definição de campo
func (d *SchemaFieldDefinition) Validate(fieldName string, value any) error {
	switch d.Type {
	case FieldTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("field '%s' must be string", fieldName)
		}
	case FieldTypeNumber:
		switch value.(type) {
		case float64, float32, int, int64, int32:
			// ok
		default:
			return fmt.Errorf("field '%s' must be number", fieldName)
		}
	case FieldTypeBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("field '%s' must be boolean", fieldName)
		}
	case FieldTypeArray:
		arr, ok := value.([]any)
		if !ok {
			return fmt.Errorf("field '%s' must be array", fieldName)
		}
		// Validar elementos do array se Items estiver definido
		if d.Items != nil {
			for _, elem := range arr {
				if err := d.Items.Validate(fieldName, elem); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
