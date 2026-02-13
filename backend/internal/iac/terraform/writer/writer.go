package writer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
	"github.com/zclconf/go-cty/cty"
)

// RenderMainTF renders a Terraform main.tf from a slice of Terraform blocks.
// Each block is written to the top-level body of the HCL file, sorting attributes for stable output.
// Returns the content of the main.tf as a string, or an error if rendering fails.
func RenderMainTF(blocks []mapper.TerraformBlock) (string, error) {
	f := hclwrite.NewEmptyFile()
	root := f.Body()

	for _, b := range blocks {
		if b.Kind == "" {
			return "", fmt.Errorf("block kind is empty")
		}
		blk := root.AppendNewBlock(b.Kind, b.Labels)
		body := blk.Body()

		// Collect and sort attribute keys for stable output.
		keys := make([]string, 0, len(b.Attributes))
		for k := range b.Attributes {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Render each attribute value using tokensForTerraformValue.
		for _, k := range keys {
			v := b.Attributes[k]
			toks, err := tokensForTerraformValue(v)
			if err != nil {
				return "", fmt.Errorf("block %s %v: attribute %q: %w", b.Kind, b.Labels, k, err)
			}
			body.SetAttributeRaw(k, toks)
		}

		// Render nested blocks (e.g., ingress, egress for security groups)
		if err := renderNestedBlocks(body, b.NestedBlocks); err != nil {
			return "", fmt.Errorf("block %s %v: %w", b.Kind, b.Labels, err)
		}

		root.AppendNewline()
	}

	return string(f.Bytes()), nil
}

// renderNestedBlocks recursively renders nested blocks
func renderNestedBlocks(body *hclwrite.Body, blocks map[string][]mapper.NestedBlock) error {
	if len(blocks) == 0 {
		return nil
	}

	// Sort nested block types for stable output
	nestedTypes := make([]string, 0, len(blocks))
	for k := range blocks {
		nestedTypes = append(nestedTypes, k)
	}
	sort.Strings(nestedTypes)

	for _, nestedType := range nestedTypes {
		for _, nested := range blocks[nestedType] {
			nestedBlk := body.AppendNewBlock(nestedType, nil)
			nestedBody := nestedBlk.Body()

			// Render attributes
			keys := make([]string, 0, len(nested.Attributes))
			for k := range nested.Attributes {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := nested.Attributes[k]
				toks, err := tokensForTerraformValue(v)
				if err != nil {
					return fmt.Errorf("nested block %s attribute %q: %w", nestedType, k, err)
				}
				nestedBody.SetAttributeRaw(k, toks)
			}

			// Recursively render nested blocks
			if err := renderNestedBlocks(nestedBody, nested.NestedBlocks); err != nil {
				return err
			}
		}
	}
	return nil
}

// RenderVariablesTF renders a variables.tf as a string from the provided variable definitions.
// Returns an empty string if no vars are provided, or an error if rendering fails.
func RenderVariablesTF(vars []mapper.Variable) (string, error) {
	if len(vars) == 0 {
		return "", nil
	}

	f := hclwrite.NewEmptyFile()
	root := f.Body()

	for _, v := range vars {
		if v.Name == "" {
			return "", fmt.Errorf("variable name is empty")
		}
		blk := root.AppendNewBlock("variable", []string{v.Name})
		body := blk.Body()

		// Optionally set variable description.
		if v.Description != "" {
			desc := v.Description
			body.SetAttributeValue("description", cty.StringVal(desc))
		}
		// Set the variable type as a raw expression (not quoted).
		if v.Type != "" {
			body.SetAttributeRaw("type", tokensForExpr(mapper.TerraformExpr(v.Type)))
		}
		// Optionally set default value.
		if v.Default != nil {
			toks, err := tokensForTerraformValue(*v.Default)
			if err != nil {
				return "", err
			}
			body.SetAttributeRaw("default", toks)
		}
		// Optionally set sensitive flag.
		if v.Sensitive {
			body.SetAttributeValue("sensitive", cty.BoolVal(true))
		}

		root.AppendNewline()
	}

	return string(f.Bytes()), nil
}

// RenderOutputsTF renders an outputs.tf as a string from the provided output definitions.
// Returns an empty string if no outputs are provided, or an error if rendering fails.
func RenderOutputsTF(outputs []mapper.Output) (string, error) {
	if len(outputs) == 0 {
		return "", nil
	}

	f := hclwrite.NewEmptyFile()
	root := f.Body()

	for _, o := range outputs {
		if o.Name == "" {
			return "", fmt.Errorf("output name is empty")
		}
		blk := root.AppendNewBlock("output", []string{o.Name})
		body := blk.Body()

		// Set the output value (required).
		toks, err := tokensForTerraformValue(o.Value)
		if err != nil {
			return "", fmt.Errorf("output %q value: %w", o.Name, err)
		}
		body.SetAttributeRaw("value", toks)

		// Optionally set output description.
		if o.Description != "" {
			body.SetAttributeValue("description", cty.StringVal(o.Description))
		}

		// Optionally set sensitive flag.
		if o.Sensitive {
			body.SetAttributeValue("sensitive", cty.BoolVal(true))
		}

		root.AppendNewline()
	}

	return string(f.Bytes()), nil
}

// tokensForTerraformValue converts a TerraformValue to hclwrite tokens suitable for setting as an attribute value.
// Supports Terraform expressions, strings, numbers, booleans, maps, and lists.
func tokensForTerraformValue(v mapper.TerraformValue) (hclwrite.Tokens, error) {
	switch {
	case v.Expr != nil:
		// Reference-style HCL expressions (unquoted).
		return tokensForExpr(*v.Expr), nil
	case v.String != nil:
		// Plain string value.
		return hclwrite.TokensForValue(cty.StringVal(*v.String)), nil
	case v.Number != nil:
		// Numeric value (float64).
		return hclwrite.TokensForValue(cty.NumberFloatVal(*v.Number)), nil
	case v.Bool != nil:
		// Boolean value.
		return hclwrite.TokensForValue(cty.BoolVal(*v.Bool)), nil
	case v.Map != nil:
		// Map values.
		return tokensForMap(v.Map)
	case v.List != nil:
		// List values.
		return tokensForList(v.List)
	default:
		// Null or unset values are not representable, error out explicitly.
		return nil, fmt.Errorf("empty TerraformValue")
	}
}

// tokensForExpr renders a TerraformExpr as HCL tokens, using attribute traversal.
// If the expression can't be split, treats it as a literal string.
func tokensForExpr(expr mapper.TerraformExpr) hclwrite.Tokens {
	parts := strings.Split(string(expr), ".")
	if len(parts) == 0 || parts[0] == "" {
		return hclwrite.TokensForValue(cty.StringVal(string(expr)))
	}

	trav := hcl.Traversal{
		hcl.TraverseRoot{Name: parts[0]},
	}
	for _, p := range parts[1:] {
		if p == "" {
			continue
		}
		trav = append(trav, hcl.TraverseAttr{Name: p})
	}
	return hclwrite.TokensForTraversal(trav)
}

// tokensForList converts a slice of TerraformValue items into HCL list tokens.
// Each list entry is rendered recursively via tokensForTerraformValue.
func tokensForList(list []mapper.TerraformValue) (hclwrite.Tokens, error) {
	toks := hclwrite.Tokens{
		&hclwrite.Token{Type: hclsyntax.TokenOBrack, Bytes: []byte("[")},
	}
	for i, item := range list {
		itemToks, err := tokensForTerraformValue(item)
		if err != nil {
			return nil, err
		}
		toks = append(toks, itemToks...)
		if i != len(list)-1 {
			toks = append(toks, &hclwrite.Token{Type: hclsyntax.TokenComma, Bytes: []byte(",")})
			toks = append(toks, &hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")})
		}
	}
	toks = append(toks, &hclwrite.Token{Type: hclsyntax.TokenCBrack, Bytes: []byte("]")})
	return toks, nil
}

// tokensForMap converts a string-keyed map of TerraformValues into HCL map tokens.
// Attribute keys are sorted for stable output.
// Each value is rendered recursively via tokensForTerraformValue.
func tokensForMap(m map[string]mapper.TerraformValue) (hclwrite.Tokens, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	toks := hclwrite.Tokens{
		&hclwrite.Token{Type: hclsyntax.TokenOBrace, Bytes: []byte("{")},
		&hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")},
	}

	for _, k := range keys {
		v := m[k]
		// Render key and = separator.
		toks = append(toks, &hclwrite.Token{Type: hclsyntax.TokenIdent, Bytes: []byte(k)})
		toks = append(toks, &hclwrite.Token{Type: hclsyntax.TokenEqual, Bytes: []byte(" = ")})
		valToks, err := tokensForTerraformValue(v)
		if err != nil {
			return nil, err
		}
		toks = append(toks, valToks...)
		toks = append(toks, &hclwrite.Token{Type: hclsyntax.TokenNewline, Bytes: []byte("\n")})
	}

	toks = append(toks, &hclwrite.Token{Type: hclsyntax.TokenCBrace, Bytes: []byte("}")})
	return toks, nil
}
