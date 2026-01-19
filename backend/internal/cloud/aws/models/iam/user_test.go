package iam

import (
	"testing"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid-user",
			user: User{
				Name: "test-user",
			},
			wantErr: false,
		},
		{
			name: "missing-name",
			user: User{},
			wantErr: true,
			errMsg:  "user name is required",
		},
		{
			name: "name-invalid-characters",
			user: User{
				Name: "test@user#invalid",
			},
			wantErr: true,
			errMsg:  "user name contains invalid characters",
		},
		{
			name: "invalid-permissions-boundary",
			user: User{
				Name:                "test-user",
				PermissionsBoundary: stringPtr("invalid-arn"),
			},
			wantErr: true,
			errMsg:  "permissions boundary must be a valid IAM policy ARN",
		},
		{
			name: "valid-user-with-permissions-boundary",
			user: User{
				Name:                "test-user",
				PermissionsBoundary: stringPtr("arn:aws:iam::123456789012:policy/boundary-policy"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
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

