package common

import (
	"fmt"
	"strings"
)

// SupportedRegions contains a list of commonly used AWS regions
var SupportedRegions = []string{
	"us-east-1",      // N. Virginia
	"us-east-2",      // Ohio
	"us-west-1",      // N. California
	"us-west-2",      // Oregon
	"eu-west-1",      // Ireland
	"eu-west-2",      // London
	"eu-central-1",   // Frankfurt
	"ap-southeast-1", // Singapore
	"ap-southeast-2", // Sydney
	"ap-northeast-1", // Tokyo
	"sa-east-1",      // São Paulo
}

// DefaultRegion is the default region to use
const DefaultRegion = "us-east-1"

// ValidateRegion checks if a region is valid
func ValidateRegion(region string) error {
	if region == "" {
		return fmt.Errorf("region cannot be empty")
	}

	for _, supported := range SupportedRegions {
		if supported == region {
			return nil
		}
	}

	return fmt.Errorf("unsupported region: %s. Supported regions: %v", region, SupportedRegions)
}

// GetRegionOrDefault returns the provided region if valid, otherwise returns default
func GetRegionOrDefault(region string) string {
	if err := ValidateRegion(region); err == nil {
		return region
	}
	return DefaultRegion
}

// ListRegions returns all supported regions
func ListRegions() []string {
	return SupportedRegions
}

// FormatRegionName formats a region code into a human-readable name
func FormatRegionName(region string) string {
	regionNames := map[string]string{
		"us-east-1":      "US East (N. Virginia)",
		"us-east-2":      "US East (Ohio)",
		"us-west-1":      "US West (N. California)",
		"us-west-2":      "US West (Oregon)",
		"eu-west-1":      "Europe (Ireland)",
		"eu-west-2":      "Europe (London)",
		"eu-central-1":   "Europe (Frankfurt)",
		"ap-southeast-1": "Asia Pacific (Singapore)",
		"ap-southeast-2": "Asia Pacific (Sydney)",
		"ap-northeast-1": "Asia Pacific (Tokyo)",
		"sa-east-1":      "South America (São Paulo)",
	}

	if name, ok := regionNames[region]; ok {
		return fmt.Sprintf("%s (%s)", name, region)
	}
	return region
}

// GetRegionCode extracts region code from formatted name
func GetRegionCode(formattedName string) string {
	parts := strings.Split(formattedName, "(")
	if len(parts) > 1 {
		code := strings.TrimSuffix(strings.TrimSpace(parts[1]), ")")
		return code
	}
	return formattedName
}

// SelectRegion displays available regions and allows selection
// For use cases, this is a simple function that returns the default or provided region
func SelectRegion(preferredRegion string) string {
	if preferredRegion != "" {
		if err := ValidateRegion(preferredRegion); err == nil {
			return preferredRegion
		}
		fmt.Printf("Warning: Invalid region '%s', using default: %s\n", preferredRegion, DefaultRegion)
	}
	return DefaultRegion
}

// DisplayRegions prints all available regions
func DisplayRegions() {
	fmt.Println("Available AWS Regions:")
	for i, region := range SupportedRegions {
		fmt.Printf("  %d. %s\n", i+1, FormatRegionName(region))
	}
}
