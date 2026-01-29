package mapper

import (
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

type fakeMapper struct {
	provider string
}

func (f *fakeMapper) Provider() string { return f.provider }

func (f *fakeMapper) SupportsResource(resourceType string) bool { return true }

func (f *fakeMapper) MapResource(res *resource.Resource) ([]TerraformBlock, error) {
	return nil, nil
}

func TestMapperRegistry_RegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	m := &fakeMapper{provider: "aws"}

	if err := reg.Register(m); err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}

	got, ok := reg.Get("aws")
	if !ok {
		t.Fatalf("Get() returned !ok for registered mapper")
	}
	if got != m {
		t.Fatalf("Get() = %v, want %v", got, m)
	}
}

func TestMapperRegistry_RegisterNil(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(nil); err == nil {
		t.Fatalf("Register(nil) expected error, got nil")
	}
}

func TestMapperRegistry_RegisterEmptyProvider(t *testing.T) {
	reg := NewRegistry()
	m := &fakeMapper{provider: ""}
	if err := reg.Register(m); err == nil {
		t.Fatalf("Register(mapper with empty provider) expected error, got nil")
	}
}

func TestMapperRegistry_RegisterDuplicate(t *testing.T) {
	reg := NewRegistry()
	m1 := &fakeMapper{provider: "aws"}
	m2 := &fakeMapper{provider: "aws"}

	if err := reg.Register(m1); err != nil {
		t.Fatalf("first Register() error = %v, want nil", err)
	}
	if err := reg.Register(m2); err == nil {
		t.Fatalf("second Register() expected error for duplicate provider, got nil")
	}
}

func TestMapperRegistry_MustGetPanicsOnMissing(t *testing.T) {
	reg := NewRegistry()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("MustGet() expected to panic for missing mapper, but did not")
		}
	}()

	_ = reg.MustGet("missing")
}

