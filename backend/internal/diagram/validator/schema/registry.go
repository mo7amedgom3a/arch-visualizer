package schema

import (
	"fmt"
	"sync"
)

// SchemaRegistry manages resource schemas for validation
type SchemaRegistry interface {
	// Register adds a schema to the registry
	Register(schema *ResourceSchema) error

	// Get retrieves a schema by resource type and provider
	Get(resourceType, provider string) (*ResourceSchema, bool)

	// GetAll returns all schemas for a provider
	GetAll(provider string) []*ResourceSchema

	// Has checks if a schema exists
	Has(resourceType, provider string) bool
}

// InMemorySchemaRegistry is an in-memory implementation of SchemaRegistry
type InMemorySchemaRegistry struct {
	mu      sync.RWMutex
	schemas map[string]map[string]*ResourceSchema // provider -> resourceType -> schema
}

// NewSchemaRegistry creates a new in-memory schema registry
func NewSchemaRegistry() *InMemorySchemaRegistry {
	return &InMemorySchemaRegistry{
		schemas: make(map[string]map[string]*ResourceSchema),
	}
}

// Register adds a schema to the registry
func (r *InMemorySchemaRegistry) Register(schema *ResourceSchema) error {
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}
	if schema.ResourceType == "" {
		return fmt.Errorf("resource type is required")
	}
	if schema.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.schemas[schema.Provider]; !exists {
		r.schemas[schema.Provider] = make(map[string]*ResourceSchema)
	}

	r.schemas[schema.Provider][schema.ResourceType] = schema
	return nil
}

// Get retrieves a schema by resource type and provider
func (r *InMemorySchemaRegistry) Get(resourceType, provider string) (*ResourceSchema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if providerSchemas, exists := r.schemas[provider]; exists {
		if schema, exists := providerSchemas[resourceType]; exists {
			return schema, true
		}
	}
	return nil, false
}

// GetAll returns all schemas for a provider
func (r *InMemorySchemaRegistry) GetAll(provider string) []*ResourceSchema {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var schemas []*ResourceSchema
	if providerSchemas, exists := r.schemas[provider]; exists {
		for _, schema := range providerSchemas {
			schemas = append(schemas, schema)
		}
	}
	return schemas
}

// Has checks if a schema exists
func (r *InMemorySchemaRegistry) Has(resourceType, provider string) bool {
	_, exists := r.Get(resourceType, provider)
	return exists
}

// DefaultRegistry is the global default schema registry
var DefaultRegistry = NewSchemaRegistry()

// Register registers a schema to the default registry
func Register(schema *ResourceSchema) error {
	return DefaultRegistry.Register(schema)
}

// Get retrieves a schema from the default registry
func Get(resourceType, provider string) (*ResourceSchema, bool) {
	return DefaultRegistry.Get(resourceType, provider)
}
