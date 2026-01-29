package terraform

import (
	"context"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
)

type fakeEngine struct {
	lastArch    *architecture.Architecture
	lastSorted  []*resource.Resource
	output      *iac.Output
	returnError error
}

func (f *fakeEngine) Name() string { return "terraform" }

func (f *fakeEngine) Generate(ctx context.Context, arch *architecture.Architecture, sorted []*resource.Resource) (*iac.Output, error) {
	f.lastArch = arch
	f.lastSorted = sorted
	return f.output, f.returnError
}

func TestCompiler_Compile_ErrorsOnNilEngine(t *testing.T) {
	c := NewCompiler(nil)
	_, err := c.Compile(context.Background(), &architecture.Architecture{})
	if err == nil {
		t.Fatalf("Compile() with nil engine expected error, got nil")
	}
}

func TestCompiler_Compile_ErrorsOnNilArchitecture(t *testing.T) {
	engine := &fakeEngine{}
	c := NewCompiler(engine)

	_, err := c.Compile(context.Background(), nil)
	if err == nil {
		t.Fatalf("Compile() with nil architecture expected error, got nil")
	}
}

func TestCompiler_Compile_HappyPath(t *testing.T) {
	engine := &fakeEngine{
		output: &iac.Output{
			Files: []iac.GeneratedFile{
				{Path: "main.tf", Content: "test", Type: "hcl"},
			},
		},
	}
	c := NewCompiler(engine)

	res := &resource.Resource{
		ID:   "r1",
		Name: "res1",
		Type: resource.ResourceType{Name: "VPC"},
	}
	arch := &architecture.Architecture{
		Resources: []*resource.Resource{res},
	}

	out, err := c.Compile(context.Background(), arch)
	if err != nil {
		t.Fatalf("Compile() error = %v, want nil", err)
	}
	if out != engine.output {
		t.Fatalf("Compile() output = %#v, want %#v", out, engine.output)
	}

	if engine.lastArch != arch {
		t.Fatalf("engine.Generate arch arg = %#v, want %#v", engine.lastArch, arch)
	}
	if len(engine.lastSorted) != 1 || engine.lastSorted[0] != res {
		t.Fatalf("engine.Generate sorted arg = %#v, want [%#v]", engine.lastSorted, res)
	}
}

