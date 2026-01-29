package mapper

import (
	"fmt"
	"sync"
)

// MapperRegistry provides provider-based selection of Terraform mappers.
type MapperRegistry struct {
	mu      sync.RWMutex
	mappers map[string]ResourceMapper // provider -> mapper
}

func NewRegistry() *MapperRegistry {
	return &MapperRegistry{
		mappers: make(map[string]ResourceMapper),
	}
}

func (r *MapperRegistry) Register(m ResourceMapper) error {
	if m == nil {
		return fmt.Errorf("mapper is nil")
	}
	provider := m.Provider()
	if provider == "" {
		return fmt.Errorf("mapper provider is empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mappers[provider]; exists {
		return fmt.Errorf("mapper for provider %q already registered", provider)
	}
	r.mappers[provider] = m
	return nil
}

func (r *MapperRegistry) Get(provider string) (ResourceMapper, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.mappers[provider]
	return m, ok
}

func (r *MapperRegistry) MustGet(provider string) ResourceMapper {
	m, ok := r.Get(provider)
	if !ok {
		panic(fmt.Sprintf("terraform mapper for provider %q not registered", provider))
	}
	return m
}

