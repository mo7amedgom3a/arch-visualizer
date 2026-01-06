package networking

import (
	"errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// NATGateway represents an AWS-specific NAT Gateway
type NATGateway struct {
	Name          string `json:"name"`
	SubnetID      string `json:"subnet_id"`
	AllocationID  string `json:"allocation_id"` // Elastic IP allocation ID
	Tags          []configs.Tag `json:"tags"`
}

// Validate performs AWS-specific validation
func (ngw *NATGateway) Validate() error {
	if ngw.Name == "" {
		return errors.New("nat gateway name is required")
	}
	if ngw.SubnetID == "" {
		return errors.New("nat gateway subnet_id is required")
	}
	if ngw.AllocationID == "" {
		return errors.New("nat gateway allocation_id is required")
	}
	return nil
}
