package iam

import (
	"testing"
)

func TestRole_Validate(t *testing.T) {
	validAssumeRolePolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}`

	tests := []struct {
		name    string
		role    Role
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid-role",
			role: Role{
				Name:             "test-role",
				AssumeRolePolicy: validAssumeRolePolicy,
			},
			wantErr: false,
		},
		{
			name: "missing-name",
			role: Role{
				AssumeRolePolicy: validAssumeRolePolicy,
			},
			wantErr: true,
			errMsg:  "role name is required",
		},
		{
			name: "missing-assume-role-policy",
			role: Role{
				Name: "test-role",
			},
			wantErr: true,
			errMsg:  "assume role policy is required",
		},
		{
			name: "invalid-assume-role-policy-json",
			role: Role{
				Name:             "test-role",
				AssumeRolePolicy: `{invalid json}`,
			},
			wantErr: true,
			errMsg:  "assume role policy must be valid JSON",
		},
		{
			name: "invalid-managed-policy-arn",
			role: Role{
				Name:             "test-role",
				AssumeRolePolicy: validAssumeRolePolicy,
				ManagedPolicyARNs: []string{"invalid-arn"},
			},
			wantErr: true,
			errMsg:  "managed policy ARN",
		},
		{
			name: "valid-role-with-managed-policies",
			role: Role{
				Name:             "test-role",
				AssumeRolePolicy: validAssumeRolePolicy,
				ManagedPolicyARNs: []string{"arn:aws:iam::123456789012:policy/test-policy"},
			},
			wantErr: false,
		},
		{
			name: "invalid-inline-policy-json",
			role: Role{
				Name:             "test-role",
				AssumeRolePolicy: validAssumeRolePolicy,
				InlinePolicies: []InlinePolicy{
					{Name: "test-inline", Policy: `{invalid json}`},
				},
			},
			wantErr: true,
			errMsg:  "inline policy",
		},
		{
			name: "valid-role-with-inline-policies",
			role: Role{
				Name:             "test-role",
				AssumeRolePolicy: validAssumeRolePolicy,
				InlinePolicies: []InlinePolicy{
					{Name: "test-inline", Policy: `{"Version":"2012-10-17","Statement":[]}`},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.role.Validate()
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

