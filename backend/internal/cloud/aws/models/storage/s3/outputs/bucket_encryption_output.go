package outputs

// BucketEncryptionOutput represents response for bucket encryption operations
type BucketEncryptionOutput struct {
	ID                 string  `json:"id"` // bucket name
	BucketKeyEnabled   bool    `json:"bucket_key_enabled"`
	SSEAlgorithm       string  `json:"sse_algorithm"`
	KMSMasterKeyID     *string `json:"kms_master_key_id,omitempty"`
}
