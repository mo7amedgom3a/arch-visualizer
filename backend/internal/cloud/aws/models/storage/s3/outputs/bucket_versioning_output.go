package outputs

// BucketVersioningOutput represents response for bucket versioning operations
type BucketVersioningOutput struct {
	ID        string  `json:"id"` // bucket name
	Status    string  `json:"status"`
	MFADelete *string `json:"mfa_delete,omitempty"`
}
