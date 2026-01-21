package storage

import "fmt"

// S3BucketVersioning models versioning configuration for a bucket
type S3BucketVersioning struct {
	Bucket string
	Status string // Enabled | Suspended | Disabled
	MFADelete *string // Enabled | Disabled (optional)
}

// Validate ensures versioning configuration is valid
func (v *S3BucketVersioning) Validate() error {
	if v.Bucket == "" {
		return fmt.Errorf("bucket is required")
	}

	validStatus := map[string]bool{
		"Enabled":   true,
		"Suspended": true,
		"Disabled":  true,
	}

	if !validStatus[v.Status] {
		return fmt.Errorf("invalid status: %s", v.Status)
	}

	if v.MFADelete != nil {
		if *v.MFADelete != "Enabled" && *v.MFADelete != "Disabled" {
			return fmt.Errorf("invalid mfa_delete: %s", *v.MFADelete)
		}
	}

	return nil
}
