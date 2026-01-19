package iam

import (
	"testing"
)

func TestPolicy_Validate(t *testing.T) {
	tests := []struct {
		name    string
		policy  Policy
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid-policy",
			policy: Policy{
				Name:          "test-policy",
				PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`,
			},
			wantErr: false,
		},
		{
			name: "missing-name",
			policy: Policy{
				PolicyDocument: `{"Version":"2012-10-17","Statement":[]}`,
			},
			wantErr: true,
			errMsg:  "policy name is required",
		},
		{
			name: "missing-policy-document",
			policy: Policy{
				Name: "test-policy",
			},
			wantErr: true,
			errMsg:  "policy document is required",
		},
		{
			name: "invalid-policy-document-json",
			policy: Policy{
				Name:          "test-policy",
				PolicyDocument: `{invalid json}`,
			},
			wantErr: true,
			errMsg:  "policy document must be valid JSON",
		},
		{
			name: "name-invalid-characters",
			policy: Policy{
				Name:          "test@policy#invalid",
				PolicyDocument: `{"Version":"2012-10-17","Statement":[]}`,
			},
			wantErr: true,
			errMsg:  "policy name contains invalid characters",
		},
		{
			name: "path-invalid-format",
			policy: Policy{
				Name:          "test-policy",
				PolicyDocument: `{"Version":"2012-10-17","Statement":[]}`,
				Path:          stringPtr("invalid-path"),
			},
			wantErr: true,
			errMsg:  "policy path must start with '/'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

