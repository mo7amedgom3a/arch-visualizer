package networking

import "errors"

// InternetGateway represents a cloud-agnostic internet gateway
type InternetGateway struct {
	ID    string
	ARN   *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
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
