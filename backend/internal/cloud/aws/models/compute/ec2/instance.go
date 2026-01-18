package ec2

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Instance represents an AWS EC2 instance configuration
type Instance struct {
	Name                     string        `json:"name"`
	AMI                      string        `json:"ami"`
	InstanceType             string        `json:"instance_type"`
	SubnetID                 string        `json:"subnet_id"`
	VpcSecurityGroupIds      []string      `json:"vpc_security_group_ids"`
	AssociatePublicIPAddress *bool         `json:"associate_public_ip_address,omitempty"`
	KeyName                  *string       `json:"key_name,omitempty"`
	IAMInstanceProfile       *string       `json:"iam_instance_profile,omitempty"`
	UserData                 *string       `json:"user_data,omitempty"`
	RootVolumeID             *string       `json:"root_volume_id,omitempty"` // Reference to storage volume
	Tags                     []configs.Tag `json:"tags,omitempty"`
}

// Validate performs AWS-specific validation
func (i *Instance) Validate() error {
	// Basic required fields
	if i.Name == "" {
		return errors.New("instance name is required")
	}
	if i.AMI == "" {
		return errors.New("ami is required")
	}
	if i.InstanceType == "" {
		return errors.New("instance type is required")
	}
	if i.SubnetID == "" {
		return errors.New("subnet id is required")
	}
	if len(i.VpcSecurityGroupIds) == 0 {
		return errors.New("at least one security group is required")
	}

	// AMI format validation
	if !strings.HasPrefix(i.AMI, "ami-") {
		return errors.New("ami must start with 'ami-'")
	}
	if len(i.AMI) < 12 || len(i.AMI) > 21 {
		return errors.New("ami format is invalid")
	}

	// Instance type validation (basic format check)
	if !IsValidInstanceType(i.InstanceType) {
		return fmt.Errorf("invalid instance type format: %s", i.InstanceType)
	}

	// Subnet ID format validation
	if !strings.HasPrefix(i.SubnetID, "subnet-") {
		return errors.New("subnet id must start with 'subnet-'")
	}

	// Security group ID format validation
	for _, sgID := range i.VpcSecurityGroupIds {
		if !strings.HasPrefix(sgID, "sg-") {
			return fmt.Errorf("security group id must start with 'sg-': %s", sgID)
		}
	}

	// UserData validation (16KB limit when base64 encoded)
	if i.UserData != nil {
		// Base64 encoding increases size by ~33%, so we check raw size
		// 16KB base64 ≈ 12KB raw
		if len(*i.UserData) > 12288 {
			return errors.New("user data cannot exceed 12KB (16KB when base64 encoded)")
		}
	}

	// Root volume ID validation
	if i.RootVolumeID != nil && *i.RootVolumeID != "" {
		if !strings.HasPrefix(*i.RootVolumeID, "vol-") {
			return errors.New("root volume id must start with 'vol-'")
		}
	}

	return nil
}

// IsValidInstanceType performs basic instance type format validation
// Valid formats: t3.micro, m5.large, c5.xlarge, etc.
// Exported so it can be used by other packages (e.g., launch_template)
func IsValidInstanceType(instanceType string) bool {
	if len(instanceType) < 3 {
		return false
	}

	parts := strings.Split(instanceType, ".")
	if len(parts) != 2 {
		return false
	}

	// Family should be alphanumeric (e.g., t3, m5, c5)
	family := parts[0]
	if len(family) < 1 || len(family) > 10 {
		return false
	}

	// Size should be a valid size (nano, micro, small, medium, large, xlarge, 2xlarge, etc.)
	size := parts[1]
	validSizes := map[string]bool{
		"nano":     true,
		"micro":    true,
		"small":    true,
		"medium":   true,
		"large":    true,
		"xlarge":   true,
		"2xlarge":  true,
		"4xlarge":  true,
		"8xlarge":  true,
		"12xlarge": true,
		"16xlarge": true,
		"24xlarge": true,
		"32xlarge": true,
		"48xlarge": true,
		"metal":    true,
	}

	return validSizes[size]
}
