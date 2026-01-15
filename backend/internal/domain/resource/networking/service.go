package networking

import "context"

// NetworkingService defines the interface for networking resource operations
// This is cloud-agnostic and can be implemented by any cloud provider
type NetworkingService interface {
	// VPC operations
	CreateVPC(ctx context.Context, vpc *VPC) (*VPC, error)
	GetVPC(ctx context.Context, id string) (*VPC, error)
	UpdateVPC(ctx context.Context, vpc *VPC) (*VPC, error)
	DeleteVPC(ctx context.Context, id string) error
	ListVPCs(ctx context.Context, region string) ([]*VPC, error)
	
	// Subnet operations
	CreateSubnet(ctx context.Context, subnet *Subnet) (*Subnet, error)
	GetSubnet(ctx context.Context, id string) (*Subnet, error)
	UpdateSubnet(ctx context.Context, subnet *Subnet) (*Subnet, error)
	DeleteSubnet(ctx context.Context, id string) error
	ListSubnets(ctx context.Context, vpcID string) ([]*Subnet, error)
	
	// Internet Gateway operations
	CreateInternetGateway(ctx context.Context, igw *InternetGateway) (*InternetGateway, error)
	AttachInternetGateway(ctx context.Context, igwID, vpcID string) error
	DetachInternetGateway(ctx context.Context, igwID, vpcID string) error
	DeleteInternetGateway(ctx context.Context, id string) error
	
	// Route Table operations
	CreateRouteTable(ctx context.Context, rt *RouteTable) (*RouteTable, error)
	GetRouteTable(ctx context.Context, id string) (*RouteTable, error)
	AssociateRouteTable(ctx context.Context, rtID, subnetID string) error
	DisassociateRouteTable(ctx context.Context, associationID string) error
	DeleteRouteTable(ctx context.Context, id string) error
	
	// Security Group operations
	CreateSecurityGroup(ctx context.Context, sg *SecurityGroup) (*SecurityGroup, error)
	GetSecurityGroup(ctx context.Context, id string) (*SecurityGroup, error)
	UpdateSecurityGroup(ctx context.Context, sg *SecurityGroup) (*SecurityGroup, error)
	DeleteSecurityGroup(ctx context.Context, id string) error
	
	// NAT Gateway operations
	CreateNATGateway(ctx context.Context, ngw *NATGateway) (*NATGateway, error)
	GetNATGateway(ctx context.Context, id string) (*NATGateway, error)
	DeleteNATGateway(ctx context.Context, id string) error
	
	// Elastic IP operations
	AllocateElasticIP(ctx context.Context, eip *ElasticIP) (*ElasticIP, error)
	GetElasticIP(ctx context.Context, id string) (*ElasticIP, error)
	ReleaseElasticIP(ctx context.Context, id string) error
	AssociateElasticIP(ctx context.Context, allocationID, instanceID string) error
	DisassociateElasticIP(ctx context.Context, associationID string) error
	ListElasticIPs(ctx context.Context, region string) ([]*ElasticIP, error)
	
	// Network ACL operations
	CreateNetworkACL(ctx context.Context, acl *NetworkACL) (*NetworkACL, error)
	GetNetworkACL(ctx context.Context, id string) (*NetworkACL, error)
	DeleteNetworkACL(ctx context.Context, id string) error
	AddNetworkACLRule(ctx context.Context, aclID string, rule ACLRule) error
	RemoveNetworkACLRule(ctx context.Context, aclID string, ruleNumber int, ruleType ACLRuleType) error
	AssociateNetworkACLWithSubnet(ctx context.Context, aclID, subnetID string) error
	DisassociateNetworkACLFromSubnet(ctx context.Context, associationID string) error
	ListNetworkACLs(ctx context.Context, vpcID string) ([]*NetworkACL, error)
	
	// Network Interface operations
	CreateNetworkInterface(ctx context.Context, eni *NetworkInterface) (*NetworkInterface, error)
	GetNetworkInterface(ctx context.Context, id string) (*NetworkInterface, error)
	DeleteNetworkInterface(ctx context.Context, id string) error
	AttachNetworkInterface(ctx context.Context, eniID, instanceID string, deviceIndex int) error
	DetachNetworkInterface(ctx context.Context, attachmentID string) error
	AssignPrivateIPAddress(ctx context.Context, eniID, privateIP string) error
	UnassignPrivateIPAddress(ctx context.Context, eniID, privateIP string) error
	ListNetworkInterfaces(ctx context.Context, subnetID string) ([]*NetworkInterface, error)
}

// NetworkingRepository defines the interface for networking resource persistence
// This abstracts data access and can be implemented for different storage backends
type NetworkingRepository interface {
	// VPC persistence
	SaveVPC(ctx context.Context, vpc *VPC) error
	FindVPCByID(ctx context.Context, id string) (*VPC, error)
	FindVPCsByRegion(ctx context.Context, region string) ([]*VPC, error)
	DeleteVPC(ctx context.Context, id string) error
	
	// Subnet persistence
	SaveSubnet(ctx context.Context, subnet *Subnet) error
	FindSubnetByID(ctx context.Context, id string) (*Subnet, error)
	FindSubnetsByVPC(ctx context.Context, vpcID string) ([]*Subnet, error)
	DeleteSubnet(ctx context.Context, id string) error
	
	// Internet Gateway persistence
	SaveInternetGateway(ctx context.Context, igw *InternetGateway) error
	FindInternetGatewayByID(ctx context.Context, id string) (*InternetGateway, error)
	DeleteInternetGateway(ctx context.Context, id string) error
	
	// Route Table persistence
	SaveRouteTable(ctx context.Context, rt *RouteTable) error
	FindRouteTableByID(ctx context.Context, id string) (*RouteTable, error)
	DeleteRouteTable(ctx context.Context, id string) error
	
	// Security Group persistence
	SaveSecurityGroup(ctx context.Context, sg *SecurityGroup) error
	FindSecurityGroupByID(ctx context.Context, id string) (*SecurityGroup, error)
	DeleteSecurityGroup(ctx context.Context, id string) error
	
	// NAT Gateway persistence
	SaveNATGateway(ctx context.Context, ngw *NATGateway) error
	FindNATGatewayByID(ctx context.Context, id string) (*NATGateway, error)
	DeleteNATGateway(ctx context.Context, id string) error
}
