package s3

import "fmt"

// BucketVersioning represents AWS S3 bucket versioning configuration
type BucketVersioning struct {
	Bucket string `json:"bucket"`
	Status string `json:"status"` // Enabled | Suspended | Disabled
	MFADelete *string `json:"mfa_delete,omitempty"` // Enabled | Disabled
}

// Validate validates versioning configuration
func (v *BucketVersioning) Validate() error {
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
