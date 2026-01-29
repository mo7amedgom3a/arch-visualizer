package generator

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

type fakeAWSMapper struct {
	mapped []*resource.Resource
}

func (m *fakeAWSMapper) Provider() string { return "aws" }

func (m *fakeAWSMapper) SupportsResource(resourceType string) bool { return true }

func (m *fakeAWSMapper) MapResource(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	m.mapped = append(m.mapped, res)
	return []tfmapper.TerraformBlock{
		{
			Kind:   "resource",
			Labels: []string{"aws_test", "example"},
		},
	}, nil
}

func TestEngine_Name(t *testing.T) {
	e := NewEngine(tfmapper.NewRegistry())
	if got := e.Name(); got != "terraform" {
		t.Fatalf("Name() = %q, want %q", got, "terraform")
	}
}

func TestEngine_Generate_ErrorsOnNilArchitecture(t *testing.T) {
	e := NewEngine(tfmapper.NewRegistry())
	_, err := e.Generate(context.Background(), nil, nil)
	if err == nil {
		t.Fatalf("Generate() with nil architecture expected error, got nil")
	}
}

func TestEngine_Generate_ErrorsOnMissingMapper(t *testing.T) {
	reg := tfmapper.NewRegistry()
	e := NewEngine(reg)

	arch := &architecture.Architecture{
		Provider: resource.AWS,
	}

	_, err := e.Generate(context.Background(), arch, []*resource.Resource{})
	if err == nil {
		t.Fatalf("Generate() expected error when no mapper is registered, got nil")
	}
}

func TestEngine_Generate_HappyPath(t *testing.T) {
	reg := tfmapper.NewRegistry()
	mapper := &fakeAWSMapper{}
	if err := reg.Register(mapper); err != nil {
		t.Fatalf("Register mapper error = %v, want nil", err)
	}

	e := NewEngine(reg)

	res := &resource.Resource{
		ID:   "vpc-1",
		Name: "main-vpc",
		Type: resource.ResourceType{
			Name: "VPC",
		},
		Provider: resource.AWS,
		Region:   "us-east-1",
		
	}

	arch := &architecture.Architecture{
		Resources: []*resource.Resource{res},
		Region:    "us-east-1",
		Provider:  resource.AWS,
	}

	out, err := e.Generate(context.Background(), arch, []*resource.Resource{res})
	fmt.Println("out", out)
	if err != nil {
		t.Fatalf("Generate() error = %v, want nil", err)
	}
	if out == nil {
		t.Fatalf("Generate() returned nil output")
	}
	if len(out.Files) != 1 {
		t.Fatalf("expected 1 generated file, got %d", len(out.Files))
	}
	main := out.Files[0]
	if main.Path != "main.tf" {
		t.Fatalf("expected main.tf path, got %q", main.Path)
	}
	if main.Type != "hcl" {
		t.Fatalf("expected file type hcl, got %q", main.Type)
	}
	if !strings.Contains(main.Content, `provider "aws"`) {
		t.Fatalf("expected provider block in content, got:\n%s", main.Content)
	}
	if len(mapper.mapped) != 1 || mapper.mapped[0] != res {
		t.Fatalf("expected mapper to receive resource, got %#v", mapper.mapped)
	}
}

