package iam

import (
	"testing"
)

func TestInstanceProfile_Validate(t *testing.T) {
	tests := []struct {
		name    string
		profile *InstanceProfile
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid profile with name",
			profile: &InstanceProfile{
				Name: "test-profile",
				Path: stringPtr("/"),
			},
			wantErr: false,
		},
		{
			name: "valid profile with name prefix",
			profile: &InstanceProfile{
				NamePrefix: stringPtr("test-"),
				Path:       stringPtr("/"),
			},
			wantErr: false,
		},
		{
			name: "missing name and name prefix",
			profile: &InstanceProfile{
				Path: stringPtr("/"),
			},
			wantErr: true,
			errMsg:  "either name or name_prefix must be provided",
		},
		{
			name: "both name and name prefix provided",
			profile: &InstanceProfile{
				Name:      "test-profile",
				NamePrefix: stringPtr("test-"),
				Path:      stringPtr("/"),
			},
			wantErr: true,
			errMsg:  "name and name_prefix cannot both be provided",
		},
		{
			name: "invalid name characters",
			profile: &InstanceProfile{
				Name: "test profile", // space is invalid
				Path: stringPtr("/"),
			},
			wantErr: true,
			errMsg:  "instance profile name contains invalid characters",
		},
		{
			name: "name too long",
			profile: &InstanceProfile{
				Name: string(make([]byte, 129)), // 129 characters
				Path: stringPtr("/"),
			},
			wantErr: true,
			errMsg:  "instance profile name must be between 1 and 128 characters",
		},
		{
			name: "invalid path - no leading slash",
			profile: &InstanceProfile{
				Name: "test-profile",
				Path: stringPtr("invalid"),
			},
			wantErr: true,
			errMsg:  "instance profile path must start with '/'",
		},
		{
			name: "valid profile with role",
			profile: &InstanceProfile{
				Name: "test-profile",
				Path: stringPtr("/"),
				Role: stringPtr("test-role"),
			},
			wantErr: false,
		},
		{
			name: "invalid role name",
			profile: &InstanceProfile{
				Name: "test-profile",
				Path: stringPtr("/"),
				Role: stringPtr("invalid role"), // space is invalid
			},
			wantErr: true,
			errMsg:  "role name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.profile.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceProfile.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("InstanceProfile.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}
