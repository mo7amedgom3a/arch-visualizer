package iam

import (
	"testing"
)

func TestGroup_Validate(t *testing.T) {
	tests := []struct {
		name    string
		group   Group
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid-group",
			group: Group{
				Name: "test-group",
			},
			wantErr: false,
		},
		{
			name: "missing-name",
			group: Group{},
			wantErr: true,
			errMsg:  "group name is required",
		},
		{
			name: "name-invalid-characters",
			group: Group{
				Name: "test@group#invalid",
			},
			wantErr: true,
			errMsg:  "group name contains invalid characters",
		},
		{
			name: "path-invalid-format",
			group: Group{
				Name: "test-group",
				Path: stringPtr("invalid-path"),
			},
			wantErr: true,
			errMsg:  "group path must start with '/'",
		},
		{
			name: "valid-group-with-path",
			group: Group{
				Name: "test-group",
				Path: stringPtr("/team-a/"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.Validate()
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

