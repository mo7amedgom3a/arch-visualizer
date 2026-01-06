package networking

import (
	"errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// InternetGateway represents an AWS-specific Internet Gateway
type InternetGateway struct {
	Name   string `json:"name"`
	VPCID  string `json:"vpc_id"`
	Tags   []configs.Tag `json:"tags"`
}

// Validate performs AWS-specific validation
func (igw *InternetGateway) Validate() error {
	if igw.Name == "" {
		return errors.New("internet gateway name is required")
	}
	if igw.VPCID == "" {
		return errors.New("internet gateway vpc_id is required")
	}
	return nil
}
