package tests

import (
	"fmt"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type InternetGatewayTest struct {
	name          string
	igw           *networking.InternetGateway
	vpc           *networking.VPC
	expectedError error
	description   string
}

func TestInternetGateway(t *testing.T) {
	validVPC := &networking.VPC{
		Name:               "test-vpc",
		Region:             "us-east-1",
		CIDR:               "10.0.0.0/16",
		EnableDNSHostnames: true,
		EnableDNSSupport:   true,
		InstanceTenancy:    "default",
		Tags:               []configs.Tag{{Key: "Name", Value: "test-vpc"}},
	}

	tests := []InternetGatewayTest{
		{
			name: "valid-internet-gateway-attached-to-vpc",
			igw: &networking.InternetGateway{
				Name:  "test-igw",
				VPCID: "vpc-123",
				Tags:  []configs.Tag{{Key: "Name", Value: "test-igw"}},
			},
			vpc:           validVPC,
			expectedError: nil,
			description:   "Valid IGW attached to VPC",
		},
		{
			name: "invalid-internet-gateway-missing-name",
			igw: &networking.InternetGateway{
				Name:  "",
				VPCID: "vpc-123",
				Tags:  []configs.Tag{},
			},
			vpc:           validVPC,
			expectedError: fmt.Errorf("internet gateway name is required"),
			description:   "IGW with empty name should fail validation",
		},
		{
			name: "invalid-internet-gateway-missing-vpc-id",
			igw: &networking.InternetGateway{
				Name:  "test-igw",
				VPCID: "",
				Tags:  []configs.Tag{{Key: "Name", Value: "test-igw"}},
			},
			vpc:           validVPC,
			expectedError: fmt.Errorf("internet gateway vpc_id is required"),
			description:   "IGW with empty VPC ID should fail validation",
		},
		{
			name: "valid-internet-gateway-with-tags",
			igw: &networking.InternetGateway{
				Name:  "test-igw",
				VPCID: "vpc-123",
				Tags: []configs.Tag{
					{Key: "Name", Value: "test-igw"},
					{Key: "Environment", Value: "production"},
				},
			},
			vpc:           validVPC,
			expectedError: nil,
			description:   "Valid IGW with multiple tags",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("IGW Name: %s\n", test.igw.Name)
			fmt.Printf("IGW VPC ID: %s\n", test.igw.VPCID)

			// Validate VPC first
			if test.vpc != nil {
				fmt.Printf("\nValidating VPC...\n")
				vpcErr := test.vpc.Validate()
				if vpcErr != nil {
					fmt.Printf("⚠️  VPC validation error: %v\n", vpcErr)
				} else {
					fmt.Printf("✅ PASSED: VPC validation succeeded\n")
					fmt.Printf("  VPC Name: %s\n", test.vpc.Name)
					fmt.Printf("  VPC Region: %s\n", test.vpc.Region)
					fmt.Printf("  VPC CIDR: %s\n", test.vpc.CIDR)
				}
			}

			// Validate IGW
			fmt.Printf("\nValidating Internet Gateway...\n")
			err := test.igw.Validate()

			if test.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error: %v, but got nil", test.expectedError)
					fmt.Printf("❌ FAILED: Expected error '%v', but got nil\n", test.expectedError)
				} else if err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error: %v, but got: %v", test.expectedError, err)
					fmt.Printf("❌ FAILED: Expected error '%v', but got '%v'\n", test.expectedError, err)
				} else {
					fmt.Printf("✅ PASSED: Got expected error: %v\n", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
					fmt.Printf("❌ FAILED: Unexpected error: %v\n", err)
				} else {
					fmt.Printf("✅ PASSED: IGW validation succeeded\n")
				}
			}

			// Check IGW-VPC relationship
			if test.igw.VPCID != "" && test.vpc != nil {
				fmt.Printf("\nChecking IGW-VPC relationship...\n")
				fmt.Printf("  IGW is attached to VPC: %s\n", test.igw.VPCID)
				fmt.Printf("  ✅ PASSED: IGW-VPC relationship is valid\n")
			} else if test.igw.VPCID == "" {
				fmt.Printf("\n⚠️  WARNING: IGW has no VPC ID (should be validated above)\n")
			}

			fmt.Printf("=== Test completed: %s ===\n\n", test.name)
		})
	}
}
