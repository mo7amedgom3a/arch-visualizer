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

// RegionAZs maps regions to their availability zones
var RegionAZs = map[string][]string{
	"us-east-1":      {"us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d", "us-east-1e", "us-east-1f"},
	"us-east-2":      {"us-east-2a", "us-east-2b", "us-east-2c"},
	"us-west-1":      {"us-west-1a", "us-west-1c"},
	"us-west-2":      {"us-west-2a", "us-west-2b", "us-west-2c", "us-west-2d"},
	"eu-west-1":      {"eu-west-1a", "eu-west-1b", "eu-west-1c"},
	"eu-west-2":      {"eu-west-2a", "eu-west-2b", "eu-west-2c"},
	"eu-central-1":   {"eu-central-1a", "eu-central-1b", "eu-central-1c"},
	"ap-southeast-1": {"ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"},
	"ap-southeast-2": {"ap-southeast-2a", "ap-southeast-2b", "ap-southeast-2c"},
	"ap-northeast-1": {"ap-northeast-1a", "ap-northeast-1c", "ap-northeast-1d"},
	"sa-east-1":      {"sa-east-1a", "sa-east-1b", "sa-east-1c"},
}

// GetAZsForRegion returns availability zones for a given region
func GetAZsForRegion(region string) []string {
	if azs, ok := RegionAZs[region]; ok {
		return azs
	}
	// Return a default mock list if not found, or empty
	return []string{region + "a", region + "b", region + "c"}
}
