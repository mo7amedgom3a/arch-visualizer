package database

import (
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/database"
	"github.com/stretchr/testify/assert"
)

func TestMapRDSInstance(t *testing.T) {
	tests := []struct {
		name     string
		input    *database.RDSInstance
		expected map[string]interface{} // Simplified expectation check
		wantErr  bool
	}{
		{
			name: "Basic RDS Instance",
			input: &database.RDSInstance{
				Name:              "mydb",
				Engine:            "postgres",
				EngineVersion:     "14.1",
				InstanceClass:     "db.t3.micro",
				AllocatedStorage:  20,
				StorageType:       "gp2",
				Username:          "admin",
				Password:          "securepassword",
				SkipFinalSnapshot: true,
			},
			expected: map[string]interface{}{
				"identifier":          "mydb",
				"engine":              "postgres",
				"engine_version":      "14.1",
				"instance_class":      "db.t3.micro",
				"allocated_storage":   20.0, // float64 in TerraformValue
				"username":            "admin",
				"password":            "securepassword",
				"skip_final_snapshot": true,
			},
			wantErr: false,
		},
		{
			name: "RDS with Tags and Security Groups",
			input: &database.RDSInstance{
				Name:                "mydb-complete",
				Engine:              "mysql",
				EngineVersion:       "8.0",
				InstanceClass:       "db.m5.large",
				AllocatedStorage:    100,
				VpcSecurityGroupIds: []string{"sg-123456", "sg-789012"},
				Tags: []configs.Tag{
					{Key: "Environment", Value: "Production"},
				},
			},
			expected: map[string]interface{}{
				"identifier": "mydb-complete",
				"engine":     "mysql",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapRDSInstance(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, "resource", got.Kind)
			assert.Equal(t, "aws_db_instance", got.Labels[0])
			assert.Equal(t, tt.input.Name, got.Labels[1])

			// Helper to check attributes
			for k, v := range tt.expected {
				attr, ok := got.Attributes[k]
				assert.True(t, ok, "attribute %s missing", k)

				// Simplified value checking
				if str, ok := v.(string); ok {
					assert.Equal(t, str, *attr.String)
				} else if num, ok := v.(float64); ok {
					assert.Equal(t, num, *attr.Number)
				} else if b, ok := v.(bool); ok {
					assert.Equal(t, b, *attr.Bool)
				}
			}

			// Special check for arrays/maps if needed
			if len(tt.input.VpcSecurityGroupIds) > 0 {
				sgAttr, ok := got.Attributes["vpc_security_group_ids"]
				assert.True(t, ok)
				assert.Equal(t, len(tt.input.VpcSecurityGroupIds), len(sgAttr.List))
			}
		})
	}
}
