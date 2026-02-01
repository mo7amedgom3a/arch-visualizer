package pricing_importer

// EC2PricingData represents pricing data for an EC2 instance
type EC2PricingData struct {
	OnDemand     string              `json:"ondemand,omitempty"`
	Reserved     *map[string]string   `json:"reserved,omitempty"`
	SpotMin      *string              `json:"spot_min,omitempty"`
	SpotMax      *string              `json:"spot_max,omitempty"`
	EMR          string              `json:"emr,omitempty"`
	PCTInterrupt string              `json:"pct_interrupt,omitempty"`
	PCTSavingsOD *int                `json:"pct_savings_od,omitempty"`
	SpotAvg      string              `json:"spot_avg,omitempty"`
}

// EC2Instance represents an EC2 instance from the scraper output
type EC2Instance struct {
	InstanceType string                          `json:"instance_type"`
	Pricing      map[string]map[string]interface{} `json:"pricing"` // map[region]map[os]EC2PricingData
}

// ImportStats tracks import statistics
type ImportStats struct {
	TotalInstances   int
	TotalRates       int
	RegionsProcessed map[string]int
	OSProcessed      map[string]int
	Errors           []string
}
