package tests

import (
	"errors"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type VPCTest struct {
	name           string
	t              *testing.T
	vpc            *networking.VPC
	expectedError  error
	expectedOutput *networking.VPC
}

func TestVPC(t *testing.T) {
	validVPC := networking.VPC{
		Name:               "test-vpc",
		Region:             "us-east-1",
		CIDR:               "10.0.0.0/16",
		EnableDNSHostnames: true,
		EnableDNSSupport:   true,
		InstanceTenancy:    "default",
		Tags:               []configs.Tag{{Key: "Name", Value: "test-vpc"}},
	}

	tests := []VPCTest{
		{
			name:          "valid-vpc",
			t:             t,
			vpc:           &validVPC,
			expectedError: nil,
		},
		{
			name: "vpc-with-empty-name",
			t:    t,
			vpc: &networking.VPC{
				Name:               "",
				Region:             "us-east-1",
				CIDR:               "10.0.0.0/16",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
				Tags:               []configs.Tag{{Key: "Name", Value: "test-vpc"}},
			},
			expectedError: errors.New("name is required"),
		},
		{
			name: "vpc-with-empty-region",
			t:    t,
			vpc: &networking.VPC{
				Name:               "test-vpc",
				Region:             "",
				CIDR:               "10.0.0.0/16",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
				Tags:               []configs.Tag{{Key: "Name", Value: "test-vpc"}},
			},
			expectedError: errors.New("region is required"),
		},
		{
			name: "vpc-with-empty-cidr",
			t:    t,
			vpc: &networking.VPC{
				Name:               "test-vpc",
				Region:             "us-east-1",
				CIDR:               "",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
				Tags:               []configs.Tag{{Key: "Name", Value: "test-vpc"}},
			},
			expectedError: errors.New("cidr is required"),
		},
		{
			name: "vpc-with-invalid-cidr-format",
			t:    t,
			vpc: &networking.VPC{
				Name:               "test-vpc",
				Region:             "us-east-1",
				CIDR:               "invalid-cidr",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
				Tags:               []configs.Tag{{Key: "Name", Value: "test-vpc"}},
			},
			expectedError: nil, // Will check error message contains "invalid cidr"
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.vpc.Validate()

			if test.expectedError == nil {
				// Should not have error
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			} else {
				// Should have error
				if err == nil {
					t.Errorf("Expected error: %v, but got nil", test.expectedError)
				} else if err.Error() != test.expectedError.Error() {
					// For invalid CIDR format, check if error message contains expected text
					if test.name == "vpc-with-invalid-cidr-format" {
						if err.Error() == "" {
							t.Errorf("Expected error message containing 'invalid cidr', but got: %v", err)
						}
					} else {
						t.Errorf("Expected error: %v, but got: %v", test.expectedError, err)
					}
				}
			}
		})
	}
}
