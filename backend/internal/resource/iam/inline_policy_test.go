package iam

import (
	"testing"
)

func TestInlinePolicy_Validate(t *testing.T) {
	tests := []struct {
		name    string
		policy  InlinePolicy
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid-inline-policy",
			policy: InlinePolicy{
				Name:   "test-inline-policy",
				Policy: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`,
			},
			wantErr: false,
		},
		{
			name: "missing-name",
			policy: InlinePolicy{
				Policy: `{"Version":"2012-10-17","Statement":[]}`,
			},
			wantErr: true,
			errMsg:  "inline policy name is required",
		},
		{
			name: "missing-policy-document",
			policy: InlinePolicy{
				Name: "test-inline-policy",
			},
			wantErr: true,
			errMsg:  "inline policy document is required",
		},
		{
			name: "invalid-policy-document-json",
			policy: InlinePolicy{
				Name:   "test-inline-policy",
				Policy: `{invalid json}`,
			},
			wantErr: true,
			errMsg:  "inline policy document must be valid JSON",
		},
		{
			name: "name-invalid-characters",
			policy: InlinePolicy{
				Name:   "test@inline#policy",
				Policy: `{"Version":"2012-10-17","Statement":[]}`,
			},
			wantErr: true,
			errMsg:  "inline policy name contains invalid characters",
		},
		{
			name: "name-too-long",
			policy: InlinePolicy{
				Name:   string(make([]byte, 129)), // 129 characters
				Policy: `{"Version":"2012-10-17","Statement":[]}`,
			},
			wantErr: true,
			errMsg:  "inline policy name must be between 1 and 128 characters",
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
