package networking

// NetworkResource represents any networking resource in the domain
// This interface allows for polymorphic handling of networking resources
type NetworkResource interface {
	GetID() string
	GetName() string
	GetVPCID() string // Returns empty string if not VPC-scoped
	Validate() error
}

// VPCScopedResource represents resources that belong to a VPC
type VPCScopedResource interface {
	NetworkResource
	GetVPCID() string
}

// Ensure all networking resources implement NetworkResource
var (
	_ NetworkResource = (*VPC)(nil)
	_ NetworkResource = (*Subnet)(nil)
	_ NetworkResource = (*InternetGateway)(nil)
	_ NetworkResource = (*RouteTable)(nil)
	_ NetworkResource = (*SecurityGroup)(nil)
	_ NetworkResource = (*NATGateway)(nil)
	_ NetworkResource = (*ElasticIP)(nil)
	_ NetworkResource = (*NetworkACL)(nil)
	_ NetworkResource = (*NetworkInterface)(nil)
)

// Implement NetworkResource for VPC
func (v *VPC) GetID() string    { return v.ID }
func (v *VPC) GetName() string  { return v.Name }
func (v *VPC) GetVPCID() string { return "" } // VPC is not VPC-scoped

// Implement NetworkResource for Subnet
func (s *Subnet) GetID() string    { return s.ID }
func (s *Subnet) GetName() string  { return s.Name }
func (s *Subnet) GetVPCID() string { return s.VPCID }

// Implement NetworkResource for InternetGateway
func (igw *InternetGateway) GetID() string    { return igw.ID }
func (igw *InternetGateway) GetName() string  { return igw.Name }
func (igw *InternetGateway) GetVPCID() string { return igw.VPCID }

// Implement NetworkResource for RouteTable
func (rt *RouteTable) GetID() string    { return rt.ID }
func (rt *RouteTable) GetName() string  { return rt.Name }
func (rt *RouteTable) GetVPCID() string { return rt.VPCID }

// Implement NetworkResource for SecurityGroup
func (sg *SecurityGroup) GetID() string    { return sg.ID }
func (sg *SecurityGroup) GetName() string  { return sg.Name }
func (sg *SecurityGroup) GetVPCID() string { return sg.VPCID }

// Implement NetworkResource for NATGateway
func (ngw *NATGateway) GetID() string    { return ngw.ID }
func (ngw *NATGateway) GetName() string  { return ngw.Name }
func (ngw *NATGateway) GetVPCID() string { return "" } // NAT Gateway is subnet-scoped, not directly VPC-scoped

// Implement NetworkResource for ElasticIP
func (eip *ElasticIP) GetID() string    { return eip.ID }
func (eip *ElasticIP) GetName() string  { return "" } // Elastic IP doesn't have a name
func (eip *ElasticIP) GetVPCID() string { return "" } // Elastic IP is region-scoped, not VPC-scoped

// Implement NetworkResource for NetworkACL
func (acl *NetworkACL) GetID() string    { return acl.ID }
func (acl *NetworkACL) GetName() string  { return acl.Name }
func (acl *NetworkACL) GetVPCID() string { return acl.VPCID }

// Implement NetworkResource for NetworkInterface
func (eni *NetworkInterface) GetID() string { return eni.ID }
func (eni *NetworkInterface) GetName() string {
	if eni.Description != nil {
		return *eni.Description
	}
	return ""
}
func (eni *NetworkInterface) GetVPCID() string { return eni.VPCID }
