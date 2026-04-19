package output

import "time"

// Result represents a generic result envelope for domain operations
// This provides a consistent structure for success/data/meta handling
type Result[T any] struct {
	// Data contains the actual result data
	Data T `json:"data"`
	// Metadata contains additional information about the operation
	Metadata *Metadata `json:"metadata,omitempty"`
}

// Metadata provides additional context about an operation result
type Metadata struct {
	// Timestamp when the result was generated
	Timestamp time.Time `json:"timestamp"`
	// RequestID is an optional identifier for the request
	RequestID *string `json:"request_id,omitempty"`
	// Additional metadata key-value pairs
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// NewResult creates a new Result with the given data
func NewResult[T any](data T) *Result[T] {
	return &Result[T]{
		Data: data,
		Metadata: &Metadata{
			Timestamp: time.Now(),
		},
	}
}

// NewResultWithMetadata creates a new Result with data and metadata
func NewResultWithMetadata[T any](data T, metadata *Metadata) *Result[T] {
	if metadata == nil {
		metadata = &Metadata{
			Timestamp: time.Now(),
		}
	}
	return &Result[T]{
		Data:     data,
		Metadata: metadata,
	}
}
