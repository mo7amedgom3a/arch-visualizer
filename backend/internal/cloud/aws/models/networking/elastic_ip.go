package networking

import (
	"errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ElasticIPAddressPoolType represents the type of IP address pool for AWS
type ElasticIPAddressPoolType string

const (
	// ElasticIPPoolAmazon Amazon's pool of IPv4 addresses
	ElasticIPPoolAmazon ElasticIPAddressPoolType = "amazon"
	// ElasticIPPoolBYOIP Public IPv4 address that you bring to your AWS account with BYOIP
	ElasticIPPoolBYOIP ElasticIPAddressPoolType = "byoip"
	// ElasticIPPoolCustomerOwned Customer-owned pool of IPv4 addresses created from your on-premises network for use with an Outpost
	ElasticIPPoolCustomerOwned ElasticIPAddressPoolType = "customer_owned"
	// ElasticIPPoolIPAM Allocate using an IPv4 IPAM pool
	ElasticIPPoolIPAM ElasticIPAddressPoolType = "ipam"
)

// ElasticIP represents an AWS-specific Elastic IP address allocation request
type ElasticIP struct {
	// AllocationID is optional - if provided, uses existing EIP instead of allocating new one
	AllocationID *string `json:"allocation_id,omitempty"` // e.g., "eipalloc-12345678"
	
	// For new allocations:
	AddressPoolType   *ElasticIPAddressPoolType `json:"address_pool_type,omitempty"`   // Type of IP address pool
	AddressPoolID     *string                    `json:"address_pool_id,omitempty"`     // Pool ID (for BYOIP, customer-owned, or IPAM pools)
	NetworkBorderGroup *string                    `json:"network_border_group,omitempty"` // Network border group (AZs, Local Zones, Wavelength Zones)
	Region            string                     `json:"region"`                        // Region where EIP is allocated
	
	// Optional tags
	Tags []configs.Tag `json:"tags,omitempty"`
}

// Validate performs AWS-specific validation
func (eip *ElasticIP) Validate() error {
	// If AllocationID is provided, we're using an existing EIP
	// In this case, we don't need to validate allocation settings
	if eip.AllocationID != nil && *eip.AllocationID != "" {
		// Validate allocation ID format
		if len(*eip.AllocationID) < 9 || (*eip.AllocationID)[:9] != "eipalloc-" {
			return errors.New("invalid allocation id format, must start with 'eipalloc-'")
		}
		return nil
	}

	// For new allocations, Region is required
	if eip.Region == "" {
		return errors.New("region is required for new elastic ip allocations")
	}

	// If AddressPoolType is specified, validate it
	if eip.AddressPoolType != nil {
		validTypes := map[ElasticIPAddressPoolType]bool{
			ElasticIPPoolAmazon:       true,
			ElasticIPPoolBYOIP:        true,
			ElasticIPPoolCustomerOwned: true,
			ElasticIPPoolIPAM:         true,
		}
		if !validTypes[*eip.AddressPoolType] {
			return errors.New("invalid address pool type")
		}

		// If pool type requires a pool ID, validate it's provided
		if *eip.AddressPoolType != ElasticIPPoolAmazon {
			if eip.AddressPoolID == nil || *eip.AddressPoolID == "" {
				return errors.New("address pool id is required for non-amazon pool types")
			}
		}
	}

	return nil
}

// IsUsingExistingEIP returns true if the EIP is using an existing allocation ID
func (eip *ElasticIP) IsUsingExistingEIP() bool {
	return eip.AllocationID != nil && *eip.AllocationID != ""
}
