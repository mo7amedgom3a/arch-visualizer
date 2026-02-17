package services

// FieldDescriptor describes a single input field for an AWS resource.
type FieldDescriptor struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"` // "string", "bool", "int", "[]string", "object"
	Required    bool     `json:"required"`
	Enum        []string `json:"enum"`
	Default     any      `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
}

// ResourceSchema describes the full input/output shape of an AWS resource.
type ResourceSchema struct {
	Label   string            `json:"label"`
	Fields  []FieldDescriptor `json:"fields"`
	Outputs map[string]string `json:"outputs"` // output field name → simple type string
}
