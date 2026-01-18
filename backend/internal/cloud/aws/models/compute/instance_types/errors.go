package compute

import "errors"

var (
	// ErrInstanceTypeNameRequired is returned when instance type name is missing
	ErrInstanceTypeNameRequired = errors.New("instance type name is required")

	// ErrInvalidCategory is returned when an invalid category is provided
	ErrInvalidCategory = errors.New("invalid instance category")

	// ErrInvalidVCPU is returned when VCPU count is invalid
	ErrInvalidVCPU = errors.New("VCPU count must be greater than 0")

	// ErrInvalidMemory is returned when memory is invalid
	ErrInvalidMemory = errors.New("memory must be greater than 0")

	// ErrInstanceTypeNotFound is returned when an instance type is not found
	ErrInstanceTypeNotFound = errors.New("instance type not found")
)
