package terraform

import (
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/stretchr/testify/assert"
)

func TestMapVPCEndpoint_Gateway(t *testing.T) {
	// Setup
	vpcID := "vpc-0a1b2c3d4e5f6g7h8"
	serviceName := "com.amazonaws.us-east-1.s3"

	res := &resource.Resource{
		Name:     "s3-endpoint",
		Type:     resource.ResourceType{Name: "VPCEndpoint"},
		Provider: "aws",
		Metadata: map[string]interface{}{
			"vpc_id":            vpcID,
			"service_name":      serviceName,
			"vpc_endpoint_type": "Gateway",
			"route_table_ids":   []interface{}{"rtb-0a1b2c3d4e5f6g7h8"},
			"policy":            "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":\"*\",\"Action\":\"*\",\"Resource\":\"*\"}]}",
		},
	}

	// Execute
	blocks, err := MapVPCEndpoint(res)

	// Verify
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)

	block := blocks[0]
	assert.Equal(t, "resource", block.Kind)
	assert.Equal(t, []string{"aws_vpc_endpoint", "s3_endpoint"}, block.Labels)

	// Check Endpoint Type
	epType := block.Attributes["vpc_endpoint_type"]
	assert.NotNil(t, epType.String)
	assert.Equal(t, "Gateway", *epType.String)

	// Check VPC ID reference
	vpcAttr := block.Attributes["vpc_id"]
	assert.NotNil(t, vpcAttr.Expr)
	assert.Contains(t, string(*vpcAttr.Expr), "aws_vpc")

	// Check Route Tables
	rtIDs := block.Attributes["route_table_ids"]
	assert.NotEmpty(t, rtIDs.List)
}

func TestMapVPCEndpoint_Interface(t *testing.T) {
	// Setup
	vpcID := "vpc-0a1b2c3d4e5f6g7h8"
	serviceName := "com.amazonaws.us-east-1.ec2" // Example interface endpoint

	res := &resource.Resource{
		Name:     "ec2-endpoint",
		Type:     resource.ResourceType{Name: "VPCEndpoint"},
		Provider: "aws",
		Metadata: map[string]interface{}{
			"vpc_id":              vpcID,
			"service_name":        serviceName,
			"vpc_endpoint_type":   "Interface",
			"subnet_ids":          []interface{}{"subnet-0a1b2c3d4e5f6g7h8"},
			"security_group_ids":  []interface{}{"sg-0a1b2c3d4e5f6g7h8"},
			"private_dns_enabled": true,
		},
	}

	// Execute
	blocks, err := MapVPCEndpoint(res)

	// Verify
	assert.NoError(t, err)
	assert.Len(t, blocks, 1)

	block := blocks[0]
	assert.Equal(t, "resource", block.Kind)

	// Check Endpoint Type
	epType := block.Attributes["vpc_endpoint_type"]
	assert.NotNil(t, epType.String)
	assert.Equal(t, "Interface", *epType.String)

	// Check Subnets
	subnets := block.Attributes["subnet_ids"]
	assert.NotEmpty(t, subnets.List)

	// Check Security Groups
	sgs := block.Attributes["security_group_ids"]
	assert.NotEmpty(t, sgs.List)

	// Check Private DNS
	dns := block.Attributes["private_dns_enabled"]
	assert.NotNil(t, dns.Bool)
	assert.Equal(t, true, *dns.Bool)
}

func TestMapVPCEndpoint_MissingRequiredFields(t *testing.T) {
	// Missing VPC ID logic test
	res := &resource.Resource{
		Name:     "invalid-endpoint",
		Type:     resource.ResourceType{Name: "VPCEndpoint"},
		Provider: "aws",
		Metadata: map[string]interface{}{
			"service_name": "com.amazonaws.us-east-1.s3",
		},
	}

	_, err := MapVPCEndpoint(res)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vpc_endpoint requires vpc_id")

	// Missing Service Name logic test
	res2 := &resource.Resource{
		Name:     "invalid-endpoint-2",
		Type:     resource.ResourceType{Name: "VPCEndpoint"},
		Provider: "aws",
		Metadata: map[string]interface{}{
			"vpc_id": "vpc-123",
		},
	}

	_, err2 := MapVPCEndpoint(res2)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "missing required config \"service_name\"")
}
