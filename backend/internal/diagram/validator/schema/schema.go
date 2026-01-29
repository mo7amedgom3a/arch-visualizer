package schema

// FieldType represents the type of a config field
type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInt     FieldType = "int"
	FieldTypeFloat   FieldType = "float"
	FieldTypeBool    FieldType = "bool"
	FieldTypeCIDR    FieldType = "cidr"
	FieldTypeArray   FieldType = "array"
	FieldTypeObject  FieldType = "object"
	FieldTypeAny     FieldType = "any"
)

// FieldConstraint represents validation constraints for a field
type FieldConstraint struct {
	MinLength   *int     `json:"min_length,omitempty"`   // For strings
	MaxLength   *int     `json:"max_length,omitempty"`   // For strings
	Pattern     *string  `json:"pattern,omitempty"`      // Regex pattern
	MinValue    *float64 `json:"min_value,omitempty"`    // For numbers
	MaxValue    *float64 `json:"max_value,omitempty"`    // For numbers
	Enum        []string `json:"enum,omitempty"`         // Allowed values
	Prefix      *string  `json:"prefix,omitempty"`       // Required prefix (e.g., "ami-")
	CIDRVersion *string  `json:"cidr_version,omitempty"` // "ipv4" or "ipv6"
}

// FieldSpec defines a single config field's specification
type FieldSpec struct {
	Name        string           `json:"name"`
	Type        FieldType        `json:"type"`
	Required    bool             `json:"required"`
	Description string           `json:"description,omitempty"`
	Default     interface{}      `json:"default,omitempty"`
	Constraints *FieldConstraint `json:"constraints,omitempty"`
	// For nested objects
	NestedFields []FieldSpec `json:"nested_fields,omitempty"`
	// For arrays
	ItemType *FieldType `json:"item_type,omitempty"`
}

// ResourceSchema defines the schema for a resource type
type ResourceSchema struct {
	ResourceType string      `json:"resource_type"` // e.g., "vpc", "ec2", "subnet"
	Provider     string      `json:"provider"`      // e.g., "aws", "azure", "gcp"
	Category     string      `json:"category"`      // e.g., "networking", "compute", "storage"
	Description  string      `json:"description,omitempty"`
	Fields       []FieldSpec `json:"fields"`
	// Relationships
	ValidParentTypes []string `json:"valid_parent_types,omitempty"` // What types can contain this resource
	ValidChildTypes  []string `json:"valid_child_types,omitempty"`  // What types this resource can contain
}

// GetRequiredFields returns all required fields
func (rs *ResourceSchema) GetRequiredFields() []FieldSpec {
	var required []FieldSpec
	for _, f := range rs.Fields {
		if f.Required {
			required = append(required, f)
		}
	}
	return required
}

// GetField returns a field by name
func (rs *ResourceSchema) GetField(name string) *FieldSpec {
	for _, f := range rs.Fields {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

// HasField checks if a field exists
func (rs *ResourceSchema) HasField(name string) bool {
	return rs.GetField(name) != nil
}
