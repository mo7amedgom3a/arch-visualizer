package architecture

import (
	"testing"
	"fmt"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory" // Initialize inventory
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

func TestAWSResourceTypeMapper_MapIRTypeToResourceType(t *testing.T) {
	mapper := NewAWSResourceTypeMapper()

	tests := []struct {
		name     string
		irType   string
		wantName string
		wantErr  bool
	}{
		{
			name:     "VPC",
			irType:   "vpc",
			wantName: "VPC",
			wantErr:  false,
		},
		{
			name:     "Subnet",
			irType:   "subnet",
			wantName: "Subnet",
			wantErr:  false,
		},
		{
			name:     "EC2",
			irType:   "ec2",
			wantName: "EC2",
			wantErr:  false,
		},
		{
			name:     "Route Table",
			irType:   "route-table",
			wantName: "RouteTable",
			wantErr:  false,
		},
		{
			name:     "Security Group",
			irType:   "security-group",
			wantName: "SecurityGroup",
			wantErr:  false,
		},
		{
			name:     "Internet Gateway",
			irType:   "internet-gateway",
			wantName: "InternetGateway",
			wantErr:  false,
		},
		{
			name:     "NAT Gateway",
			irType:   "nat-gateway",
			wantName: "NATGateway",
			wantErr:  false,
		},
		{
			name:     "Elastic IP",
			irType:   "elastic-ip",
			wantName: "ElasticIP",
			wantErr:  false,
		},
		{
			name:     "Lambda",
			irType:   "lambda",
			wantName: "Lambda",
			wantErr:  false,
		},
		{
			name:     "S3",
			irType:   "s3",
			wantName: "S3",
			wantErr:  false,
		},
		{
			name:     "EBS",
			irType:   "ebs",
			wantName: "EBS",
			wantErr:  false,
		},
		{
			name:     "RDS",
			irType:   "rds",
			wantName: "RDS",
			wantErr:  false,
		},
		{
			name:     "DynamoDB",
			irType:   "dynamodb",
			wantName: "DynamoDB",
			wantErr:  false,
		},
		{
			name:     "Load Balancer",
			irType:   "load-balancer",
			wantName: "LoadBalancer",
			wantErr:  false,
		},
		{
			name:     "Auto Scaling Group",
			irType:   "auto-scaling-group",
			wantName: "AutoScalingGroup",
			wantErr:  false,
		},
		{
			name:    "Unknown IR Type",
			irType:  "unknown-type",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.MapIRTypeToResourceType(tt.irType)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapIRTypeToResourceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Fatal("Expected non-nil ResourceType")
				}

				if got.Name != tt.wantName {
					t.Errorf("MapIRTypeToResourceType() Name = %v, want %v", got.Name, tt.wantName)
				}
			}
		})
	}
}

func TestAWSResourceTypeMapper_MapResourceNameToResourceType(t *testing.T) {
	mapper := NewAWSResourceTypeMapper()

	tests := []struct {
		name         string
		resourceName string
		wantID       string
		wantCategory string
		wantKind     string
		wantErr      bool
	}{
		{
			name:         "VPC",
			resourceName: "VPC",
			wantID:       "vpc",
			wantCategory: string(resource.CategoryNetworking),
			wantKind:     "Network",
			wantErr:      false,
		},
		{
			name:         "Subnet",
			resourceName: "Subnet",
			wantID:       "subnet",
			wantCategory: string(resource.CategoryNetworking),
			wantKind:     "Network",
			wantErr:      false,
		},
		{
			name:         "EC2",
			resourceName: "EC2",
			wantID:       "ec2",
			wantCategory: string(resource.CategoryCompute),
			wantKind:     "VirtualMachine",
			wantErr:      false,
		},
		{
			name:         "RouteTable",
			resourceName: "RouteTable",
			wantID:       "route-table",
			wantCategory: string(resource.CategoryNetworking),
			wantKind:     "Network",
			wantErr:      false,
		},
		{
			name:         "SecurityGroup",
			resourceName: "SecurityGroup",
			wantID:       "security-group",
			wantCategory: string(resource.CategoryNetworking),
			wantKind:     "Network",
			wantErr:      false,
		},
		{
			name:         "InternetGateway",
			resourceName: "InternetGateway",
			wantID:       "internet-gateway",
			wantCategory: string(resource.CategoryNetworking),
			wantKind:     "Gateway",
			wantErr:      false,
		},
		{
			name:         "NATGateway",
			resourceName: "NATGateway",
			wantID:       "nat-gateway",
			wantCategory: string(resource.CategoryNetworking),
			wantKind:     "Gateway",
			wantErr:      false,
		},
		{
			name:         "ElasticIP",
			resourceName: "ElasticIP",
			wantID:       "elastic-ip",
			wantCategory: string(resource.CategoryNetworking),
			wantKind:     "Network",
			wantErr:      false,
		},
		{
			name:         "Lambda",
			resourceName: "Lambda",
			wantID:       "lambda",
			wantCategory: string(resource.CategoryCompute),
			wantKind:     "Function",
			wantErr:      false,
		},
		{
			name:         "S3",
			resourceName: "S3",
			wantID:       "s3",
			wantCategory: string(resource.CategoryStorage),
			wantKind:     "Storage",
			wantErr:      false,
		},
		{
			name:         "EBS",
			resourceName: "EBS",
			wantID:       "ebs",
			wantCategory: string(resource.CategoryStorage),
			wantKind:     "Storage",
			wantErr:      false,
		},
		{
			name:         "RDS",
			resourceName: "RDS",
			wantID:       "rds",
			wantCategory: string(resource.CategoryDatabase),
			wantKind:     "Database",
			wantErr:      false,
		},
		{
			name:         "DynamoDB",
			resourceName: "DynamoDB",
			wantID:       "dynamodb",
			wantCategory: string(resource.CategoryDatabase),
			wantKind:     "Database",
			wantErr:      false,
		},
		{
			name:         "LoadBalancer",
			resourceName: "LoadBalancer",
			wantID:       "load-balancer",
			wantCategory: string(resource.CategoryCompute),
			wantKind:     "LoadBalancer",
			wantErr:      false,
		},
		{
			name:         "AutoScalingGroup",
			resourceName: "AutoScalingGroup",
			wantID:       "auto-scaling-group",
			wantCategory: string(resource.CategoryCompute),
			wantKind:     "VirtualMachine",
			wantErr:      false,
		},
		{
			name:         "Unknown Resource Name",
			resourceName: "UnknownResource",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.MapResourceNameToResourceType(tt.resourceName)
			fmt.Println(got.Name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapResourceNameToResourceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Fatal("Expected non-nil ResourceType")
				}

				if got.ID != tt.wantID {
					t.Errorf("MapResourceNameToResourceType() ID = %v, want %v", got.ID, tt.wantID)
				}

				if got.Category != tt.wantCategory {
					t.Errorf("MapResourceNameToResourceType() Category = %v, want %v", got.Category, tt.wantCategory)
				}

				if got.Kind != tt.wantKind {
					t.Errorf("MapResourceNameToResourceType() Kind = %v, want %v", got.Kind, tt.wantKind)
				}
			}
		})
	}
}

func TestAWSResourceTypeMapper_MapResourceNameToResourceType_RegionalFlags(t *testing.T) {
	mapper := NewAWSResourceTypeMapper()

	tests := []struct {
		name         string
		resourceName string
		wantRegional bool
		wantGlobal   bool
	}{
		{
			name:         "VPC is regional",
			resourceName: "VPC",
			wantRegional: true,
			wantGlobal:   false,
		},
		{
			name:         "S3 is global",
			resourceName: "S3",
			wantRegional: false,
			wantGlobal:   true,
		},
		{
			name:         "EC2 is regional",
			resourceName: "EC2",
			wantRegional: true,
			wantGlobal:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.MapResourceNameToResourceType(tt.resourceName)
			if err != nil {
				t.Fatalf("MapResourceNameToResourceType() error = %v", err)
			}

			if got.IsRegional != tt.wantRegional {
				t.Errorf("MapResourceNameToResourceType() IsRegional = %v, want %v", got.IsRegional, tt.wantRegional)
			}

			if got.IsGlobal != tt.wantGlobal {
				t.Errorf("MapResourceNameToResourceType() IsGlobal = %v, want %v", got.IsGlobal, tt.wantGlobal)
			}
		})
	}
}

func TestAWSResourceTypeMapper_MapIRTypeToResourceType_WithAliases(t *testing.T) {
	mapper := NewAWSResourceTypeMapper()

	// Test that aliases work (these should be registered in inventory)
	// Note: This test depends on the inventory being properly initialized
	// If the inventory has aliases like "igw" for "internet-gateway", test them here

	tests := []struct {
		name     string
		irType   string
		wantName string
		wantErr  bool
	}{
		{
			name:     "Internet Gateway via alias",
			irType:   "igw", // Common alias
			wantName: "InternetGateway",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.MapIRTypeToResourceType(tt.irType)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapIRTypeToResourceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Fatal("Expected non-nil ResourceType")
				}

				if got.Name != tt.wantName {
					t.Errorf("MapIRTypeToResourceType() Name = %v, want %v", got.Name, tt.wantName)
				}
			}
		})
	}
}

func TestAWSResourceTypeMapper_Integration(t *testing.T) {
	mapper := NewAWSResourceTypeMapper()

	// Test that MapIRTypeToResourceType and MapResourceNameToResourceType
	// produce consistent results for the same resource

	irType := "vpc"
	resourceName := "VPC"

	// Map IR type to ResourceType
	rt1, err1 := mapper.MapIRTypeToResourceType(irType)
	if err1 != nil {
		t.Fatalf("MapIRTypeToResourceType() error = %v", err1)
	}

	// Map resource name to ResourceType
	rt2, err2 := mapper.MapResourceNameToResourceType(resourceName)
	if err2 != nil {
		t.Fatalf("MapResourceNameToResourceType() error = %v", err2)
	}

	// Verify they produce the same result
	if rt1.ID != rt2.ID {
		t.Errorf("ID mismatch: MapIRTypeToResourceType() = %v, MapResourceNameToResourceType() = %v", rt1.ID, rt2.ID)
	}

	if rt1.Name != rt2.Name {
		t.Errorf("Name mismatch: MapIRTypeToResourceType() = %v, MapResourceNameToResourceType() = %v", rt1.Name, rt2.Name)
	}

	if rt1.Category != rt2.Category {
		t.Errorf("Category mismatch: MapIRTypeToResourceType() = %v, MapResourceNameToResourceType() = %v", rt1.Category, rt2.Category)
	}

	if rt1.Kind != rt2.Kind {
		t.Errorf("Kind mismatch: MapIRTypeToResourceType() = %v, MapResourceNameToResourceType() = %v", rt1.Kind, rt2.Kind)
	}
}
