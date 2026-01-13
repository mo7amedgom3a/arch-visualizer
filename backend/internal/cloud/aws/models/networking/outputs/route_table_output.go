package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// RouteTableAssociation represents a route table association with a subnet
type RouteTableAssociation struct {
	ID       string `json:"id"`        // Association ID
	SubnetID string `json:"subnet_id"`
	Main     bool   `json:"main"`      // Whether this is the main route table
}

// RouteTableOutput represents AWS Route Table output/response data after creation
type RouteTableOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "rtb-12345678"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:route-table/rtb-12345678"
	
	// Configuration (from input)
	Name    string              `json:"name"`
	VPCID   string              `json:"vpc_id"`
	Routes  []networking.Route  `json:"routes"`
	
	// AWS-specific output fields
	Associations []RouteTableAssociation `json:"associations"`
	CreationTime time.Time               `json:"creation_time"`
	Tags         []configs.Tag           `json:"tags"`
}
