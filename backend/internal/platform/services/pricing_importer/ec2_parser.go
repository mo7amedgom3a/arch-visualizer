package pricing_importer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ParseEC2Instances parses the scraper JSON file and returns EC2 instances
func ParseEC2Instances(filePath string) ([]EC2Instance, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var instances []EC2Instance
	if err := json.Unmarshal(data, &instances); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return instances, nil
}

// ConvertToPricingRates converts EC2 instances to PricingRate models
func ConvertToPricingRates(instances []EC2Instance) ([]*models.PricingRate, error) {
	var rates []*models.PricingRate
	now := time.Now()

	for _, instance := range instances {
		if instance.Pricing == nil {
			continue
		}

		for region, osMap := range instance.Pricing {
			if osMap == nil {
				continue
			}

			for osName, pricingDataRaw := range osMap {
				// Handle both EC2PricingData object and string (legacy format)
				var pricingData EC2PricingData
				var onDemandStr string

				switch v := pricingDataRaw.(type) {
				case string:
					// Legacy format: direct string price
					onDemandStr = v
				case map[string]interface{}:
					// Modern format: object with ondemand field
					if ondemand, ok := v["ondemand"].(string); ok {
						onDemandStr = ondemand
					} else {
						// Try to unmarshal as EC2PricingData
						dataBytes, _ := json.Marshal(v)
						if err := json.Unmarshal(dataBytes, &pricingData); err == nil {
							onDemandStr = pricingData.OnDemand
						}
					}
				default:
					// Try to unmarshal as EC2PricingData
					dataBytes, _ := json.Marshal(v)
					if err := json.Unmarshal(dataBytes, &pricingData); err == nil {
						onDemandStr = pricingData.OnDemand
					}
				}

				if onDemandStr == "" || onDemandStr == "0" {
					continue
				}

				// Parse price string to float
				rate, err := strconv.ParseFloat(strings.TrimSpace(onDemandStr), 64)
				if err != nil {
					continue // Skip invalid prices
				}

				if rate <= 0 {
					continue // Skip zero or negative prices
				}

				// Normalize OS name
				normalizedOS := normalizeOS(osName)

				// Create PricingRate
				regionPtr := &region
				instanceTypePtr := &instance.InstanceType
				osPtr := &normalizedOS

				pricingRate := &models.PricingRate{
					Provider:        "aws",
					ResourceType:    "ec2_instance",
					ComponentName:   "EC2 Instance Hourly",
					PricingModel:    "per_hour",
					Unit:            "hour",
					Rate:            rate,
					Currency:        "USD",
					Region:          regionPtr,
					InstanceType:    instanceTypePtr,
					OperatingSystem: osPtr,
					EffectiveFrom:   now,
					EffectiveTo:     nil,
				}

				rates = append(rates, pricingRate)
			}
		}
	}

	return rates, nil
}

// normalizeOS normalizes OS names from scraper format to our format
func normalizeOS(osName string) string {
	osName = strings.ToLower(strings.TrimSpace(osName))

	// Map scraper OS names to our constants
	osMap := map[string]string{
		"linux":   "linux",
		"mswin":   "mswin",
		"windows": "mswin",
		"rhel":    "rhel",
		"suse":    "suse",
	}

	if normalized, ok := osMap[osName]; ok {
		return normalized
	}

	// Default to linux if unknown
	return "linux"
}
