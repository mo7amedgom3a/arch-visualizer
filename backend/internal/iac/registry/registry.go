package registry

import (
	"fmt"
	"sync"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
)

// EngineRegistry holds IaC engines by name (e.g. "terraform", "pulumi").
// This enables loose coupling: codegen selects an engine by name without
// importing engine-specific packages directly.
type EngineRegistry struct {
	mu      sync.RWMutex
	engines map[string]iac.Engine
}

func New() *EngineRegistry {
	return &EngineRegistry{
		engines: make(map[string]iac.Engine),
	}
}

func (r *EngineRegistry) Register(engine iac.Engine) error {
	if engine == nil {
		return fmt.Errorf("engine is nil")
	}
	name := engine.Name()
	if name == "" {
		return fmt.Errorf("engine name is empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.engines[name]; exists {
		return fmt.Errorf("engine %q already registered", name)
	}
	r.engines[name] = engine
	return nil
}

func (r *EngineRegistry) Get(name string) (iac.Engine, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.engines[name]
	return e, ok
}

func (r *EngineRegistry) MustGet(name string) iac.Engine {
	e, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprintf("iac engine %q not registered", name))
	}
	return e
}

