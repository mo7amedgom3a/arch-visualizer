package networking

import "errors"

// NATGateway represents a cloud-agnostic NAT Gateway
type NATGateway struct {
	ID              string
	Name            string
	SubnetID        string // Must be in a public subnet
	AllocationID   *string // Elastic IP allocation ID (optional)
}

// Validate performs domain-level validation
func (ngw *NATGateway) Validate() error {
	if ngw.Name == "" {
		return errors.New("nat gateway name is required")
	}
	if ngw.SubnetID == "" {
		return errors.New("nat gateway subnet_id is required")
	}
	return nil
}
