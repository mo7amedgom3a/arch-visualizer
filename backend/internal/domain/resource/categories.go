package resource

// Resource categories for cloud-agnostic classification
const (
	CategoryNetworking  = "Networking"
	CategoryCompute     = "Compute"
	CategoryStorage     = "Storage"
	CategoryDatabase    = "Database"
	CategoryContainers  = "Containers"
	CategoryIAM         = "IAM"
	CategoryMonitoring  = "Monitoring"
	CategorySecurity    = "Security"
	CategoryAnalytics   = "Analytics"
	CategoryApplication = "Application"
)

// ValidCategories returns all valid resource categories
func ValidCategories() []string {
	return []string{
		CategoryNetworking,
		CategoryCompute,
		CategoryStorage,
		CategoryDatabase,
		CategoryContainers,
		CategoryIAM,
		CategoryMonitoring,
		CategorySecurity,
		CategoryAnalytics,
		CategoryApplication,
	}
}

// IsValidCategory checks if a category string is valid
func IsValidCategory(category string) bool {
	for _, valid := range ValidCategories() {
		if valid == category {
			return true
		}
	}
	return false
}
