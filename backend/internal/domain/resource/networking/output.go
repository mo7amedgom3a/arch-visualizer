package networking

import (
	"time"
)

// VPCOutput represents the output data for a VPC after creation/update
type VPCOutput struct {
	// Core identifiers
	ID     string
	ARN    *string
	Name   string
	Region string

	// Configuration
	CIDR             string
	IPv6CIDR         *string
	EnableDNS        bool
	EnableDNSHostnames bool

	// State
	State      *string
	IsDefault  *bool
	CreatedAt  *time.Time
	OwnerID    *string
}

// SubnetOutput represents the output data for a subnet after creation/update
type SubnetOutput struct {
	// Core identifiers
	ID              string
	ARN             *string
	Name            string
	VPCID           string

	// Configuration
	CIDR            string
	AvailabilityZone *string
	IsPublic        bool

	// Output fields (cloud-generated)
	State               *string
	AvailableIPCount    *int
	MapPublicIPOnLaunch *bool

	// CreatedAt timestamp
	CreatedAt *time.Time
}

// InternetGatewayOutput represents the output data for an internet gateway after creation/update
type InternetGatewayOutput struct {
	// Core identifiers
	ID    string
	ARN   *string
	Name  string
	VPCID *string

	// State
	State           *string
	AttachmentState *string
	CreatedAt       *time.Time
}

// RouteTableOutput represents the output data for a route table after creation/update
type RouteTableOutput struct {
	// Core identifiers
	ID    string
	ARN   *string
	Name  string
	VPCID string

	// State
	State     *string
	CreatedAt *time.Time
}

// SecurityGroupOutput represents the output data for a security group after creation/update
type SecurityGroupOutput struct {
	// Core identifiers
	ID          string
	ARN         *string
	Name        string
	Description string
	VPCID       string

	// Configuration
	Rules []SecurityGroupRule

	// State
	CreatedAt *time.Time
}

// NATGatewayOutput represents the output data for a NAT gateway after creation/update
type NATGatewayOutput struct {
	// Core identifiers
	ID              string
	ARN             *string
	Name            string
	SubnetID        string
	AllocationID   *string

	// Output fields (cloud-generated)
	PublicIP *string

	// State
	State     *string
	CreatedAt *time.Time
}

// ElasticIPOutput represents the output data for an Elastic IP after allocation
type ElasticIPOutput struct {
	// Core identifiers
	ID           string
	ARN          *string
	AllocationID *string
	PublicIP     *string
	Region       string

	// Configuration
	AddressPoolType   *ElasticIPAddressPoolType
	AddressPoolID     *string
	NetworkBorderGroup *string

	// State
	State     *string
	CreatedAt *time.Time
}

// NetworkACLOutput represents the output data for a network ACL after creation/update
type NetworkACLOutput struct {
	// Core identifiers
	ID    string
	ARN   *string
	Name  string
	VPCID string

	// State
	State     *string
	CreatedAt *time.Time
}

// NetworkInterfaceOutput represents the output data for a network interface after creation/update
type NetworkInterfaceOutput struct {
	// Core identifiers
	ID       string
	ARN      *string
	Name     string
	SubnetID string

	// Configuration
	PrivateIPs      []string
	SecurityGroupIDs []string

	// Output fields (cloud-generated)
	MACAddress *string
	State      *string

	// CreatedAt timestamp
	CreatedAt *time.Time
}

// ToVPCOutput converts a VPC domain model to VPCOutput
func ToVPCOutput(vpc *VPC) *VPCOutput {
	if vpc == nil {
		return nil
	}
	return &VPCOutput{
		ID:                 vpc.ID,
		ARN:                vpc.ARN,
		Name:               vpc.Name,
		Region:             vpc.Region,
		CIDR:               vpc.CIDR,
		IPv6CIDR:           vpc.IPv6CIDR,
		EnableDNS:          vpc.EnableDNS,
		EnableDNSHostnames: vpc.EnableDNSHostnames,
	}
}

// ToSubnetOutput converts a Subnet domain model to SubnetOutput
func ToSubnetOutput(subnet *Subnet) *SubnetOutput {
	if subnet == nil {
		return nil
	}
	return &SubnetOutput{
		ID:              subnet.ID,
		ARN:             subnet.ARN,
		Name:            subnet.Name,
		VPCID:           subnet.VPCID,
		CIDR:            subnet.CIDR,
		AvailabilityZone: subnet.AvailabilityZone,
		IsPublic:        subnet.IsPublic,
	}
}

// ToSecurityGroupOutput converts a SecurityGroup domain model to SecurityGroupOutput
func ToSecurityGroupOutput(sg *SecurityGroup) *SecurityGroupOutput {
	if sg == nil {
		return nil
	}
	return &SecurityGroupOutput{
		ID:          sg.ID,
		ARN:         sg.ARN,
		Name:        sg.Name,
		Description: sg.Description,
		VPCID:       sg.VPCID,
		Rules:       sg.Rules,
	}
}

// ToNATGatewayOutput converts a NATGateway domain model to NATGatewayOutput
func ToNATGatewayOutput(ngw *NATGateway) *NATGatewayOutput {
	if ngw == nil {
		return nil
	}
	return &NATGatewayOutput{
		ID:            ngw.ID,
		ARN:           ngw.ARN,
		Name:          ngw.Name,
		SubnetID:      ngw.SubnetID,
		AllocationID:  ngw.AllocationID,
	}
}

// ToElasticIPOutput converts an ElasticIP domain model to ElasticIPOutput
func ToElasticIPOutput(eip *ElasticIP) *ElasticIPOutput {
	if eip == nil {
		return nil
	}
	return &ElasticIPOutput{
		ID:                 eip.ID,
		ARN:                eip.ARN,
		AllocationID:       eip.AllocationID,
		PublicIP:           eip.PublicIP,
		Region:             eip.Region,
		AddressPoolType:    eip.AddressPoolType,
		AddressPoolID:      eip.AddressPoolID,
		NetworkBorderGroup: eip.NetworkBorderGroup,
	}
}
