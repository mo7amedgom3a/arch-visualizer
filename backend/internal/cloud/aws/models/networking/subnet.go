package networking

import (
	"errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Subnet represents an AWS-specific subnet
type Subnet struct {
	Name             string `json:"name"`
	VPCID            string `json:"vpc_id"`
	CIDR             string `json:"cidr"`
	AvailabilityZone string `json:"availability_zone"`
	MapPublicIPOnLaunch bool `json:"map_public_ip_on_launch"`
	Tags             []configs.Tag `json:"tags"`
}

// Validate performs AWS-specific validation
func (s *Subnet) Validate() error {
	if s.Name == "" {
		return errors.New("subnet name is required")
	}
	if s.VPCID == "" {
		return errors.New("subnet vpc_id is required")
	}
	if s.CIDR == "" {
		return errors.New("subnet cidr is required")
	}
	if s.AvailabilityZone == "" {
		return errors.New("subnet availability_zone is required")
	}
	return nil
}
