package networking

import (
	"fmt"
	"net"
)

// CIDROverlaps checks if two CIDR blocks overlap
// Returns true if they overlap, false otherwise
func CIDROverlaps(cidr1, cidr2 string) (bool, error) {
	_, ipNet1, err := net.ParseCIDR(cidr1)
	if err != nil {
		return false, fmt.Errorf("invalid cidr1 format: %w", err)
	}

	_, ipNet2, err := net.ParseCIDR(cidr2)
	if err != nil {
		return false, fmt.Errorf("invalid cidr2 format: %w", err)
	}

	// Check if either network contains the other's starting IP
	return ipNet1.Contains(ipNet2.IP) || ipNet2.Contains(ipNet1.IP), nil
}

// CIDRContains checks if parentCIDR contains childCIDR
// Returns true if parentCIDR fully contains childCIDR, false otherwise
func CIDRContains(parentCIDR, childCIDR string) (bool, error) {
	_, parentNet, err := net.ParseCIDR(parentCIDR)
	if err != nil {
		return false, fmt.Errorf("invalid parent cidr format: %w", err)
	}

	_, childNet, err := net.ParseCIDR(childCIDR)
	if err != nil {
		return false, fmt.Errorf("invalid child cidr format: %w", err)
	}

	// Check if parent contains child's starting IP
	if !parentNet.Contains(childNet.IP) {
		return false, nil
	}

	// Check if child mask is more specific (larger number) than parent mask
	parentMaskSize, _ := parentNet.Mask.Size()
	childMaskSize, _ := childNet.Mask.Size()
	if childMaskSize <= parentMaskSize {
		return false, nil
	}

	// Calculate the last IP in the child network
	// The last IP is the broadcast address of the child network
	lastIP := make(net.IP, len(childNet.IP))
	copy(lastIP, childNet.IP)
	for i := range lastIP {
		lastIP[i] |= ^childNet.Mask[i]
	}

	// Check if the last IP of the child network is also within the parent
	return parentNet.Contains(lastIP), nil
}
