package tests

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
)

type InstanceTest struct {
	name          string
	instance      *ec2.Instance
	expectedError error
	description   string
}

func TestInstance(t *testing.T) {
	tests := []InstanceTest{
		{
			name: "valid-instance",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "t3.micro",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"sg-123"},
				Tags:                []configs.Tag{{Key: "Name", Value: "test-instance"}},
			},
			expectedError: nil,
			description:   "Valid EC2 instance with required fields",
		},
		{
			name: "invalid-instance-missing-name",
			instance: &ec2.Instance{
				Name:                "",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "t3.micro",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"sg-123"},
			},
			expectedError: errors.New("instance name is required"),
			description:   "Instance with empty name should fail validation",
		},
		{
			name: "invalid-instance-missing-ami",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "",
				InstanceType:        "t3.micro",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"sg-123"},
			},
			expectedError: errors.New("ami is required"),
			description:   "Instance with empty AMI should fail validation",
		},
		{
			name: "invalid-instance-missing-instance-type",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"sg-123"},
			},
			expectedError: errors.New("instance type is required"),
			description:   "Instance with empty instance type should fail validation",
		},
		{
			name: "invalid-instance-missing-subnet-id",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "t3.micro",
				SubnetID:            "",
				VpcSecurityGroupIds: []string{"sg-123"},
			},
			expectedError: errors.New("subnet id is required"),
			description:   "Instance with empty subnet ID should fail validation",
		},
		{
			name: "invalid-instance-missing-security-groups",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "t3.micro",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{},
			},
			expectedError: errors.New("at least one security group is required"),
			description:   "Instance with no security groups should fail validation",
		},
		{
			name: "invalid-instance-invalid-ami-format",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "invalid-ami",
				InstanceType:        "t3.micro",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"sg-123"},
			},
			expectedError: errors.New("ami must start with 'ami-'"),
			description:   "Instance with invalid AMI format should fail validation",
		},
		{
			name: "invalid-instance-invalid-subnet-id-format",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "t3.micro",
				SubnetID:            "invalid-subnet",
				VpcSecurityGroupIds: []string{"sg-123"},
			},
			expectedError: errors.New("subnet id must start with 'subnet-'"),
			description:   "Instance with invalid subnet ID format should fail validation",
		},
		{
			name: "invalid-instance-invalid-security-group-id-format",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "t3.micro",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"invalid-sg"},
			},
			expectedError: fmt.Errorf("security group id must start with 'sg-': %s", "invalid-sg"),
			description:   "Instance with invalid security group ID format should fail validation",
		},
		{
			name: "invalid-instance-invalid-instance-type",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "invalid-type",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"sg-123"},
			},
			expectedError: fmt.Errorf("invalid instance type format: %s", "invalid-type"),
			description:   "Instance with invalid instance type format should fail validation",
		},
		{
			name: "valid-instance-with-root-block-device",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "t3.micro",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"sg-123"},
				RootVolumeID:        stringPtr("vol-123"),
				Tags:                []configs.Tag{{Key: "Name", Value: "test-instance"}},
			},
			expectedError: nil,
			description:   "Valid instance with root block device configuration",
		},

		{
			name: "invalid-instance-user-data-too-large",
			instance: &ec2.Instance{
				Name:                "test-instance",
				AMI:                 "ami-0c55b159cbfafe1f0",
				InstanceType:        "t3.micro",
				SubnetID:            "subnet-123",
				VpcSecurityGroupIds: []string{"sg-123"},
				UserData:            stringPtr(string(make([]byte, 12289))), // 12KB + 1 byte
			},
			expectedError: errors.New("user data cannot exceed 12KB (16KB when base64 encoded)"),
			description:   "Instance with user data exceeding 12KB should fail validation",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)

			err := test.instance.Validate()

			if test.expectedError == nil {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
					fmt.Printf("❌ FAILED: Expected no error, but got: %v\n", err)
				} else {
					fmt.Printf("✅ PASSED: Instance validation succeeded\n")
				}
			} else {
				if err == nil {
					t.Errorf("Expected error: %v, but got none", test.expectedError)
					fmt.Printf("❌ FAILED: Expected error: %v, but got none\n", test.expectedError)
				} else if !strings.Contains(err.Error(), test.expectedError.Error()) {
					// For wrapped errors, check if error message contains expected text
					t.Errorf("Expected error containing: %v, but got: %v", test.expectedError, err)
					fmt.Printf("❌ FAILED: Expected error containing: %v, but got: %v\n", test.expectedError, err)
				} else {
					fmt.Printf("✅ PASSED: Validation correctly returned error: %v\n", err)
				}
			}

			fmt.Printf("=== Test completed: %s ===\n", test.name)
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
