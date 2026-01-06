package networking

import "errors"

// InternetGateway represents a cloud-agnostic internet gateway
type InternetGateway struct {
	ID    string
	Name  string
	VPCID string // Attached VPC
}

// Validate performs domain-level validation
func (igw *InternetGateway) Validate() error {
	if igw.Name == "" {
		return errors.New("internet gateway name is required")
	}
	if igw.VPCID == "" {
		return errors.New("internet gateway vpc_id is required")
	}
	return nil
}
