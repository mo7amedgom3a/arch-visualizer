package launch_template

import (
	"errors"
	"fmt"
)

// MetadataOptions represents instance metadata service options (IMDSv2)
type MetadataOptions struct {
	HTTPEndpoint            *string `json:"http_endpoint,omitempty"`             // "enabled" or "disabled"
	HTTPTokens              *string `json:"http_tokens,omitempty"`                // "required" (IMDSv2) or "optional" (IMDSv1)
	HTTPPutResponseHopLimit *int    `json:"http_put_response_hop_limit,omitempty"` // Hop limit for PUT requests (1-64)
}

// Validate performs validation on metadata options
func (mo *MetadataOptions) Validate() error {
	if mo.HTTPEndpoint != nil {
		validEndpoints := map[string]bool{
			"enabled":  true,
			"disabled": true,
		}
		if !validEndpoints[*mo.HTTPEndpoint] {
			return fmt.Errorf("invalid http_endpoint: %s (must be 'enabled' or 'disabled')", *mo.HTTPEndpoint)
		}
	}

	if mo.HTTPTokens != nil {
		validTokens := map[string]bool{
			"required": true,
			"optional": true,
		}
		if !validTokens[*mo.HTTPTokens] {
			return fmt.Errorf("invalid http_tokens: %s (must be 'required' or 'optional')", *mo.HTTPTokens)
		}
	}

	if mo.HTTPPutResponseHopLimit != nil {
		if *mo.HTTPPutResponseHopLimit < 1 || *mo.HTTPPutResponseHopLimit > 64 {
			return errors.New("http_put_response_hop_limit must be between 1 and 64")
		}
	}

	return nil
}
