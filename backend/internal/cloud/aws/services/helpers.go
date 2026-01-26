package services

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Fixed timestamp for deterministic outputs (2024-01-01 00:00:00 UTC)
var fixedTimestamp = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

// GetFixedTimestamp returns a fixed timestamp for deterministic outputs
func GetFixedTimestamp() time.Time {
	return fixedTimestamp
}

// GenerateDeterministicID generates a deterministic ID from a seed string
func GenerateDeterministicID(seed string) string {
	hash := sha256.Sum256([]byte(seed))
	return fmt.Sprintf("%x", hash[:8])
}

// GenerateInstanceID generates a deterministic instance ID from a name
func GenerateInstanceID(name string) string {
	id := GenerateDeterministicID(name)
	// Use all 16 hex characters (GenerateDeterministicID returns 16 chars from 8 bytes)
	return fmt.Sprintf("i-%s", id)
}

// GenerateARN generates a deterministic ARN for a resource
func GenerateARN(service, resourceType, resourceID, region string) string {
	accountID := "123456789012" // Fixed account ID
	return fmt.Sprintf("arn:aws:%s:%s:%s:%s/%s", service, region, accountID, resourceType, resourceID)
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// Int32Ptr returns a pointer to an int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// Int64Ptr returns a pointer to an int64
func Int64Ptr(i int64) *int64 {
	return &i
}
