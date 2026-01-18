package launch_template

import (
	"errors"
	"strings"
)

// IAMInstanceProfile represents IAM instance profile configuration for Launch Templates
type IAMInstanceProfile struct {
	Name *string `json:"name,omitempty"` // IAM instance profile name
	ARN  *string `json:"arn,omitempty"`  // IAM instance profile ARN
}

// Validate performs validation on IAM instance profile
func (iip *IAMInstanceProfile) Validate() error {
	// Either Name or ARN must be provided
	if (iip.Name == nil || *iip.Name == "") && (iip.ARN == nil || *iip.ARN == "") {
		return errors.New("IAM instance profile name or ARN is required")
	}

	// If Name is provided, validate format
	if iip.Name != nil && *iip.Name != "" {
		if len(*iip.Name) < 1 || len(*iip.Name) > 128 {
			return errors.New("IAM instance profile name must be between 1 and 128 characters")
		}
	}

	// If ARN is provided, validate format
	if iip.ARN != nil && *iip.ARN != "" {
		if !strings.HasPrefix(*iip.ARN, "arn:aws:iam::") {
			return errors.New("IAM instance profile ARN must start with 'arn:aws:iam::'")
		}
	}

	return nil
}
