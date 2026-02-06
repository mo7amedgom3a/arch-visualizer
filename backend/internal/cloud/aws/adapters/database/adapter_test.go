package database

import (
	"context"
	"testing"

	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/database"
	domaindatabase "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/database"
	"github.com/stretchr/testify/assert"
)

func TestCreateRDSInstance(t *testing.T) {
	tests := []struct {
		name        string
		input       *domaindatabase.RDSInstance
		expectedMap map[string]interface{}
		wantErr     bool
	}{
		{
			name: "Valid RDS Instance",
			input: &domaindatabase.RDSInstance{
				Name:                "my-db-1",
				Engine:              "mysql",
				EngineVersion:       "8.0",
				InstanceClass:       "db.m5.large",
				AllocatedStorage:    100,
				MultiAZ:             true,
				VpcSecurityGroupIds: []string{"sg-1", "sg-2"},
			},
			expectedMap: map[string]interface{}{
				"Engine":           "mysql",
				"AllocatedStorage": 100,
				"MultiAZ":          true,
			},
			wantErr: false,
		},
		{
			name: "Invalid Instance - Missing Name",
			input: &domaindatabase.RDSInstance{
				// Name missing
				Engine:        "mysql",
				InstanceClass: "db.t3.micro",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdapter(awsservice.NewDatabaseService())
			got, err := a.CreateRDSInstance(context.TODO(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.input.Name, got.Name)

			// Checks specific fields
			if val, ok := tt.expectedMap["Engine"]; ok {
				assert.Equal(t, val, got.Engine)
			}
			if val, ok := tt.expectedMap["AllocatedStorage"]; ok {
				assert.Equal(t, val, got.AllocatedStorage)
			}
			if val, ok := tt.expectedMap["MultiAZ"]; ok {
				assert.Equal(t, val, got.MultiAZ)
			}
		})
	}
}
