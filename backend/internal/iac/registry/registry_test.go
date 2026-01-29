package registry

import (
	"context"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
)

type fakeEngine struct {
	name string
}

func (f *fakeEngine) Name() string { return f.name }

func (f *fakeEngine) Generate(ctx context.Context, arch *architecture.Architecture, sorted []*resource.Resource) (*iac.Output, error) {
	return &iac.Output{Files: []iac.GeneratedFile{}}, nil
}

func TestEngineRegistry_RegisterAndGet(t *testing.T) {
	reg := New()

	engine := &fakeEngine{name: "terraform"}
	if err := reg.Register(engine); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	got, ok := reg.Get("terraform")
	if !ok {
		t.Fatalf("Get() returned !ok for registered engine")
	}
	if got != engine {
		t.Fatalf("Get() = %v, want %v", got, engine)
	}
}

func TestEngineRegistry_RegisterNil(t *testing.T) {
	reg := New()
	if err := reg.Register(nil); err == nil {
		t.Fatalf("Register(nil) expected error, got nil")
	}
}

func TestEngineRegistry_RegisterEmptyName(t *testing.T) {
	reg := New()
	engine := &fakeEngine{name: ""}
	if err := reg.Register(engine); err == nil {
		t.Fatalf("Register(engine with empty name) expected error, got nil")
	}
}

func TestEngineRegistry_RegisterDuplicate(t *testing.T) {
	reg := New()
	engine1 := &fakeEngine{name: "terraform"}
	engine2 := &fakeEngine{name: "terraform"}

	if err := reg.Register(engine1); err != nil {
		t.Fatalf("first Register() error = %v, want nil", err)
	}
	if err := reg.Register(engine2); err == nil {
		t.Fatalf("second Register() expected error for duplicate name, got nil")
	}
}

func TestEngineRegistry_MustGetPanicsOnMissing(t *testing.T) {
	reg := New()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("MustGet() expected to panic for missing engine, but did not")
		}
	}()

	_ = reg.MustGet("missing")
}

