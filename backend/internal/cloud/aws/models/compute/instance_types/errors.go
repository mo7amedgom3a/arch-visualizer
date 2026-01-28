package compute

import (
	awserrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/errors"
)

var (
	// ErrInstanceTypeNameRequired is returned when instance type name is missing
	// Deprecated: Use awserrors.NewInstanceTypeNameRequired() instead
	ErrInstanceTypeNameRequired = awserrors.NewInstanceTypeNameRequired()

	// ErrInvalidCategory is returned when an invalid category is provided
	// Deprecated: Use awserrors.NewInvalidCategory() instead
	ErrInvalidCategory = awserrors.NewInvalidCategory("")

	// ErrInvalidVCPU is returned when VCPU count is invalid
	// Deprecated: Use awserrors.NewInvalidVCPU() instead
	ErrInvalidVCPU = awserrors.NewInvalidVCPU(0)

	// ErrInvalidMemory is returned when memory is invalid
	// Deprecated: Use awserrors.NewInvalidMemory() instead
	ErrInvalidMemory = awserrors.NewInvalidMemory(0)

	// ErrInstanceTypeNotFound is returned when an instance type is not found
	// Deprecated: Use awserrors.NewInstanceTypeNotFound() instead
	ErrInstanceTypeNotFound = awserrors.NewInstanceTypeNotFound("")
)
