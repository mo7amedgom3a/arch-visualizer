package schema

import (
	"testing"
)

func TestSchemaRegistry(t *testing.T) {
	registry := NewSchemaRegistry()

	// Register a test schema
	testSchema := &ResourceSchema{
		ResourceType: "test-resource",
		Provider:     "test-provider",
		Category:     "test",
		Description:  "Test resource",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true},
			{Name: "count", Type: FieldTypeInt, Required: false},
		},
	}

	err := registry.Register(testSchema)
	if err != nil {
		t.Fatalf("Failed to register schema: %v", err)
	}

	// Test Get
	retrieved, exists := registry.Get("test-resource", "test-provider")
	if !exists {
		t.Fatal("Schema should exist")
	}
	if retrieved.ResourceType != "test-resource" {
		t.Errorf("Expected resource type 'test-resource', got '%s'", retrieved.ResourceType)
	}

	// Test Has
	if !registry.Has("test-resource", "test-provider") {
		t.Error("Has should return true for existing schema")
	}
	if registry.Has("non-existent", "test-provider") {
		t.Error("Has should return false for non-existent schema")
	}

	// Test GetAll
	all := registry.GetAll("test-provider")
	if len(all) != 1 {
		t.Errorf("Expected 1 schema, got %d", len(all))
	}
}

func TestResourceSchemaHelpers(t *testing.T) {
	schema := &ResourceSchema{
		ResourceType: "vpc",
		Provider:     "aws",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true},
			{Name: "cidr", Type: FieldTypeCIDR, Required: true},
			{Name: "tags", Type: FieldTypeArray, Required: false},
		},
	}

	// Test GetRequiredFields
	required := schema.GetRequiredFields()
	if len(required) != 2 {
		t.Errorf("Expected 2 required fields, got %d", len(required))
	}

	// Test GetField
	field := schema.GetField("cidr")
	if field == nil {
		t.Fatal("Field 'cidr' should exist")
	}
	if field.Type != FieldTypeCIDR {
		t.Errorf("Expected CIDR type, got %s", field.Type)
	}

	// Test HasField
	if !schema.HasField("name") {
		t.Error("Should have field 'name'")
	}
	if schema.HasField("non-existent") {
		t.Error("Should not have field 'non-existent'")
	}
}

func TestDefaultRegistryHasAWSSchemas(t *testing.T) {
	// Check that default registry has AWS schemas loaded
	awsResources := []string{
		"vpc", "subnet", "ec2", "security-group", "route-table",
		"internet-gateway", "nat-gateway", "lambda", "s3", "rds",
	}

	for _, rt := range awsResources {
		if !DefaultRegistry.Has(rt, "aws") {
			t.Errorf("Default registry should have AWS schema for '%s'", rt)
		}
	}
}

func TestAWSSchemaFields(t *testing.T) {
	// Test VPC schema
	vpcSchema, exists := DefaultRegistry.Get("vpc", "aws")
	if !exists {
		t.Fatal("VPC schema should exist")
	}

	// Check required fields
	cidrField := vpcSchema.GetField("cidr")
	if cidrField == nil || !cidrField.Required {
		t.Error("VPC should have required 'cidr' field")
	}

	nameField := vpcSchema.GetField("name")
	if nameField == nil || !nameField.Required {
		t.Error("VPC should have required 'name' field")
	}

	// Test EC2 schema
	ec2Schema, exists := DefaultRegistry.Get("ec2", "aws")
	if !exists {
		t.Fatal("EC2 schema should exist")
	}

	amiField := ec2Schema.GetField("ami")
	if amiField == nil || !amiField.Required {
		t.Error("EC2 should have required 'ami' field")
	}
	if amiField.Constraints == nil || amiField.Constraints.Prefix == nil || *amiField.Constraints.Prefix != "ami-" {
		t.Error("EC2 'ami' field should have 'ami-' prefix constraint")
	}

	// Test subnet schema
	subnetSchema, exists := DefaultRegistry.Get("subnet", "aws")
	if !exists {
		t.Fatal("Subnet schema should exist")
	}

	azField := subnetSchema.GetField("availabilityZoneId")
	if azField == nil || !azField.Required {
		t.Error("Subnet should have required 'availabilityZoneId' field")
	}

	// Check valid parent types
	if len(subnetSchema.ValidParentTypes) == 0 {
		t.Error("Subnet should have valid parent types defined")
	}
	foundVPC := false
	for _, pt := range subnetSchema.ValidParentTypes {
		if pt == "vpc" {
			foundVPC = true
			break
		}
	}
	if !foundVPC {
		t.Error("Subnet should have 'vpc' as valid parent type")
	}
}

func TestSchemaRegistration(t *testing.T) {
	registry := NewSchemaRegistry()

	// Test nil schema
	err := registry.Register(nil)
	if err == nil {
		t.Error("Should error on nil schema")
	}

	// Test empty resource type
	err = registry.Register(&ResourceSchema{Provider: "aws"})
	if err == nil {
		t.Error("Should error on empty resource type")
	}

	// Test empty provider
	err = registry.Register(&ResourceSchema{ResourceType: "test"})
	if err == nil {
		t.Error("Should error on empty provider")
	}
}
