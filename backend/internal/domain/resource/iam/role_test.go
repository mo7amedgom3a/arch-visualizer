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
			name: "name-too-long",
			role: Role{
				Name:             string(make([]byte, 65)), // 65 characters
				AssumeRolePolicy: validAssumeRolePolicy,
			},
			wantErr: true,
			errMsg:  "role name must be between 1 and 64 characters",
		},
		{
			name: "name-invalid-characters",
			role: Role{
				Name:             "test@role#invalid",
				AssumeRolePolicy: validAssumeRolePolicy,
			},
			wantErr: true,
			errMsg:  "role name contains invalid characters",
		},
		{
			name: "invalid-permissions-boundary",
			role: Role{
				Name:             "test-role",
				AssumeRolePolicy: validAssumeRolePolicy,
				PermissionsBoundary: stringPtr("invalid-arn"),
			},
			wantErr: true,
			errMsg:  "permissions boundary must be a valid IAM policy ARN",
		},
		{
			name: "valid-role-with-permissions-boundary",
			role: Role{
				Name:             "test-role",
				AssumeRolePolicy: validAssumeRolePolicy,
				PermissionsBoundary: stringPtr("arn:aws:iam::123456789012:policy/boundary-policy"),
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

