package pricing

// PricingModel represents the type of pricing model
type PricingModel string

const (
	// PerHour pricing model - charges per hour (e.g., NAT Gateway: $0.045/hour)
	PerHour PricingModel = "per_hour"
	// PerGB pricing model - charges per gigabyte (e.g., Data Transfer: $0.09/GB)
	PerGB PricingModel = "per_gb"
	// PerRequest pricing model - charges per request (e.g., API Gateway: $0.000001/request)
	PerRequest PricingModel = "per_request"
	// OneTime pricing model - one-time fees (e.g., setup fees)
	OneTime PricingModel = "one_time"
	// Tiered pricing model - tiered pricing (e.g., first 1GB free, then $0.09/GB)
	Tiered PricingModel = "tiered"
	// Percentage pricing model - percentage-based (e.g., 2% of resource cost)
	Percentage PricingModel = "percentage"
)

// Currency represents the currency type
type Currency string

const (
	// USD United States Dollar
	USD Currency = "USD"
	// EUR Euro
	EUR Currency = "EUR"
	// GBP British Pound
	GBP Currency = "GBP"
)

// Period represents the time period for cost calculations
type Period string

const (
	// Hourly period
	Hourly Period = "hourly"
	// Monthly period
	Monthly Period = "monthly"
	// Yearly period
	Yearly Period = "yearly"
)

// CloudProvider represents the cloud provider
type CloudProvider string

const (
	// AWS Amazon Web Services
	AWS CloudProvider = "aws"
	// Azure Microsoft Azure
	Azure CloudProvider = "azure"
	// GCP Google Cloud Platform
	GCP CloudProvider = "gcp"
)
