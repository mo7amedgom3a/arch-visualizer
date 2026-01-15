package networking

import "errors"

// ElasticIPAddressPoolType represents the type of IP address pool
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

// ElasticIP represents a cloud-agnostic Elastic IP address
type ElasticIP struct {
	ID                string  // Allocation ID (e.g., "eipalloc-12345678")
	ARN               *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	PublicIP          *string // Public IP address (populated after allocation)
	AllocationID      *string // Existing allocation ID (if using existing EIP)
	AddressPoolType   *ElasticIPAddressPoolType // Type of IP address pool
	AddressPoolID     *string // Pool ID (for BYOIP, customer-owned, or IPAM pools)
	NetworkBorderGroup *string // Network border group (collection of AZs, Local Zones, Wavelength Zones)
	Region            string  // Region where EIP is allocated
}

// Validate performs domain-level validation
func (eip *ElasticIP) Validate() error {
	// If AllocationID is provided, it means we're using an existing EIP
	// In this case, we don't need to validate pool settings
	if eip.AllocationID != nil && *eip.AllocationID != "" {
		return nil // Using existing EIP, no further validation needed
	}

	// For new allocations, Region is required
	if eip.Region == "" {
		return errors.New("elastic ip region is required for new allocations")
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
			return errors.New("invalid elastic ip address pool type")
		}

		// If pool type requires a pool ID, validate it's provided
		if *eip.AddressPoolType != ElasticIPPoolAmazon && eip.AddressPoolID == nil {
			return errors.New("address pool id is required for non-amazon pool types")
		}
	}

	return nil
}

// IsUsingExistingEIP returns true if the EIP is using an existing allocation ID
func (eip *ElasticIP) IsUsingExistingEIP() bool {
	return eip.AllocationID != nil && *eip.AllocationID != ""
}
