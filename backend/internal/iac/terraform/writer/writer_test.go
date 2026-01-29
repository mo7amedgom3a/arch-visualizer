package writer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func TestRenderMainTF_SimpleResourceAndProvider(t *testing.T) {
	region := "us-east-1"
	blocks := []mapper.TerraformBlock{
		{
			Kind:   "provider",
			Labels: []string{"aws"},
			Attributes: map[string]mapper.TerraformValue{
				"region": {String: &region},
			},
		},
		{
			Kind:   "resource",
			Labels: []string{"aws_vpc", "main"},
			Attributes: map[string]mapper.TerraformValue{
				"cidr_block":           {String: strPtr("10.0.0.0/16")},
				"enable_dns_hostnames": {Bool: boolPtr(true)},
				"enable_dns_support":   {Bool: boolPtr(true)},
				"instance_tenancy":     {String: strPtr("default")},
				"tags": {Map: map[string]mapper.TerraformValue{
					"Name": {String: strPtr("main")},
				}},
			},
		},
	}

	out, err := RenderMainTF(blocks)
	fmt.Println(out)
	if err != nil {
		t.Fatalf("RenderMainTF() error = %v, want nil", err)
	}

	if !strings.Contains(out, `provider "aws"`) {
		t.Fatalf("expected provider block, got:\n%s", out)
	}
	if !strings.Contains(out, `resource "aws_vpc" "main"`) {
		t.Fatalf("expected aws_vpc resource block, got:\n%s", out)
	}
	if !strings.Contains(out, "cidr_block") || !strings.Contains(out, "10.0.0.0/16") {
		t.Fatalf("expected cidr_block attribute with value, got:\n%s", out)
	}
	if !strings.Contains(out, "enable_dns_hostnames") {
		t.Fatalf("expected enable_dns_hostnames attribute, got:\n%s", out)
	}
	if !strings.Contains(out, "enable_dns_support") {
		t.Fatalf("expected enable_dns_support attribute, got:\n%s", out)
	}
	if !strings.Contains(out, "instance_tenancy") || !strings.Contains(out, "default") {
		t.Fatalf("expected instance_tenancy attribute, got:\n%s", out)
	}
	if !strings.Contains(out, "tags") || !strings.Contains(out, "Name") || !strings.Contains(out, "main") {
		t.Fatalf("expected tags attribute with Name=\"main\", got:\n%s", out)
	}
}

func TestRenderVariablesTF_TypeIsExpression(t *testing.T) {
	vars := []mapper.Variable{
		{
			Name:        "region",
			Type:        "string",
			Description: "AWS region",
		},
	}

	out, err := RenderVariablesTF(vars)
	fmt.Println(out)
	if err != nil {
		t.Fatalf("RenderVariablesTF() error = %v, want nil", err)
	}

	if !strings.Contains(out, `variable "region"`) {
		t.Fatalf("expected variable block, got:\n%s", out)
	}
	if !strings.Contains(out, "type") {
		t.Fatalf("expected type attribute, got:\n%s", out)
	}
}

func TestTokensForExpr_ProducesTraversal(t *testing.T) {
	expr := mapper.TerraformExpr("aws_vpc.main.id")
	out := tokensForExpr(expr).Bytes()

	if string(out) != "aws_vpc.main.id" {
		t.Fatalf("expected traversal aws_vpc.main.id, got %q", string(out))
	}
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }
