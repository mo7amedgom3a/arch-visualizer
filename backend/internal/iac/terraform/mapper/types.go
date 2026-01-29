package mapper

// TerraformExpr represents an HCL expression that should be written as-is
// (e.g. "aws_vpc.main.id") rather than as a quoted string literal.
type TerraformExpr string

// TerraformValue is a value that can be rendered into HCL.
// Exactly one field should be set.
type TerraformValue struct {
	String *string
	Number *float64
	Bool   *bool
	List   []TerraformValue
	Map    map[string]TerraformValue
	Expr   *TerraformExpr
}

// TerraformBlock represents a Terraform block. For resources, Kind="resource".
// For providers, Kind="provider", etc.
type TerraformBlock struct {
	Kind   string   // "resource", "provider", "data", "output", ...
	Labels []string // e.g. ["aws_vpc","main"]

	// Attributes are key -> value inside the block.
	Attributes map[string]TerraformValue
}

// Reference is a helper to express a reference like "aws_vpc.main.id".
type Reference struct {
	ResourceType string // e.g. "aws_vpc"
	ResourceName string // e.g. "main"
	Attribute    string // e.g. "id"
}

func (r Reference) Expr() TerraformExpr {
	if r.Attribute == "" {
		return TerraformExpr(r.ResourceType + "." + r.ResourceName)
	}
	return TerraformExpr(r.ResourceType + "." + r.ResourceName + "." + r.Attribute)
}

// Variable describes a Terraform input variable (optional generation).
type Variable struct {
	Name        string
	Type        string // e.g. "string", "number", "bool", "list(string)"
	Description string
	Default     *TerraformValue
	Sensitive   bool
}

