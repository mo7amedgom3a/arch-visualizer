package inventory

import (
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

func TestGetAWSResourceClassifications(t *testing.T) {
	classifications := GetAWSResourceClassifications()

	if len(classifications) == 0 {
		t.Fatal("Expected at least one resource classification")
	}

	// Verify we have resources from different categories
	categories := make(map[string]int)
	resourceNames := make(map[string]bool)

	for _, classification := range classifications {
		if classification.ResourceName == "" {
			t.Error("Resource name should not be empty")
		}

		if classification.Category == "" {
			t.Error("Category should not be empty")
		}

		if classification.IRType == "" {
			t.Error("IRType should not be empty")
		}

		// Track categories
		categories[classification.Category]++

		// Track resource names (should be unique)
		if resourceNames[classification.ResourceName] {
			t.Errorf("Duplicate resource name: %s", classification.ResourceName)
		}
		resourceNames[classification.ResourceName] = true
	}

	// Verify we have resources from expected categories
	expectedCategories := []string{
		resource.CategoryNetworking,
		resource.CategoryCompute,
		resource.CategoryStorage,
		resource.CategoryDatabase,
	}

	for _, cat := range expectedCategories {
		if categories[cat] == 0 {
			t.Errorf("Expected at least one resource in category %s", cat)
		}
	}
}

func TestGetAWSResourceClassifications_NetworkingResources(t *testing.T) {
	classifications := GetAWSResourceClassifications()

	networkingResources := []string{"VPC", "Subnet", "RouteTable", "SecurityGroup", "InternetGateway", "NATGateway", "ElasticIP"}

	found := make(map[string]bool)
	for _, classification := range classifications {
		if classification.Category == resource.CategoryNetworking {
			found[classification.ResourceName] = true

			// Verify IR type is set
			if classification.IRType == "" {
				t.Errorf("IRType should not be empty for %s", classification.ResourceName)
			}

			// Verify aliases are set
			if len(classification.Aliases) == 0 {
				t.Errorf("Aliases should not be empty for %s", classification.ResourceName)
			}
		}
	}

	for _, resName := range networkingResources {
		if !found[resName] {
			t.Errorf("Expected networking resource %s not found", resName)
		}
	}
}

func TestGetAWSResourceClassifications_ComputeResources(t *testing.T) {
	classifications := GetAWSResourceClassifications()

	computeResources := []string{"EC2", "Lambda", "LoadBalancer", "AutoScalingGroup"}

	found := make(map[string]bool)
	for _, classification := range classifications {
		if classification.Category == resource.CategoryCompute {
			found[classification.ResourceName] = true
		}
	}

	for _, resName := range computeResources {
		if !found[resName] {
			t.Errorf("Expected compute resource %s not found", resName)
		}
	}
}

func TestGetAWSResourceClassifications_StorageResources(t *testing.T) {
	classifications := GetAWSResourceClassifications()

	storageResources := []string{"S3", "EBS"}

	found := make(map[string]bool)
	for _, classification := range classifications {
		if classification.Category == resource.CategoryStorage {
			found[classification.ResourceName] = true
		}
	}

	for _, resName := range storageResources {
		if !found[resName] {
			t.Errorf("Expected storage resource %s not found", resName)
		}
	}
}

func TestGetAWSResourceClassifications_DatabaseResources(t *testing.T) {
	classifications := GetAWSResourceClassifications()

	databaseResources := []string{"RDS", "DynamoDB"}

	found := make(map[string]bool)
	for _, classification := range classifications {
		if classification.Category == resource.CategoryDatabase {
			found[classification.ResourceName] = true
		}
	}

	for _, resName := range databaseResources {
		if !found[resName] {
			t.Errorf("Expected database resource %s not found", resName)
		}
	}
}

func TestGetAWSResourceClassifications_IRTypeMapping(t *testing.T) {
	classifications := GetAWSResourceClassifications()

	// Create a test inventory and register all classifications
	inv := NewInventory()
	for _, classification := range classifications {
		inv.RegisterResource(classification, FunctionRegistry{})
	}

	// Test IR type mappings
	testCases := []struct {
		irType      string
		expectedRes string
		description string
	}{
		{"vpc", "VPC", "VPC IR type"},
		{"subnet", "Subnet", "Subnet IR type"},
		{"ec2", "EC2", "EC2 IR type"},
		{"lambda", "Lambda", "Lambda IR type"},
		{"s3", "S3", "S3 IR type"},
		{"rds", "RDS", "RDS IR type"},
		{"dynamodb", "DynamoDB", "DynamoDB IR type"},
		{"route-table", "RouteTable", "RouteTable IR type"},
		{"security-group", "SecurityGroup", "SecurityGroup IR type"},
		{"internet-gateway", "InternetGateway", "InternetGateway IR type"},
		{"nat-gateway", "NATGateway", "NATGateway IR type"},
		{"elastic-ip", "ElasticIP", "ElasticIP IR type"},
		{"load-balancer", "LoadBalancer", "LoadBalancer IR type"},
		{"auto-scaling-group", "AutoScalingGroup", "AutoScalingGroup IR type"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			resourceName, ok := inv.GetResourceNameByIRType(tc.irType)
			if !ok {
				t.Errorf("Failed to map IR type %s to resource name", tc.irType)
				return
			}

			if resourceName != tc.expectedRes {
				t.Errorf("Expected %s for IR type %s, got %s", tc.expectedRes, tc.irType, resourceName)
			}
		})
	}
}

func TestGetAWSResourceClassifications_AliasMapping(t *testing.T) {
	classifications := GetAWSResourceClassifications()

	// Create a test inventory and register all classifications
	inv := NewInventory()
	for _, classification := range classifications {
		inv.RegisterResource(classification, FunctionRegistry{})
	}

	// Test alias mappings
	testCases := []struct {
		alias       string
		expectedRes string
		description string
	}{
		{"igw", "InternetGateway", "InternetGateway alias igw"},
		{"sg", "SecurityGroup", "SecurityGroup alias sg"},
		{"eip", "ElasticIP", "ElasticIP alias eip"},
		{"nat", "NATGateway", "NATGateway alias nat"},
		{"instance", "EC2", "EC2 alias instance"},
		{"bucket", "S3", "S3 alias bucket"},
		{"alb", "LoadBalancer", "LoadBalancer alias alb"},
		{"asg", "AutoScalingGroup", "AutoScalingGroup alias asg"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			resourceName, ok := inv.GetResourceNameByIRType(tc.alias)
			if !ok {
				t.Errorf("Failed to map alias %s to resource name", tc.alias)
				return
			}

			if resourceName != tc.expectedRes {
				t.Errorf("Expected %s for alias %s, got %s", tc.expectedRes, tc.alias, resourceName)
			}
		})
	}
}
