package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRDSInstancePricing(t *testing.T) {
	tests := []struct {
		name             string
		instanceClass    string
		engine           string
		multiAZ          bool
		allocatedStorage float64
		storageType      string
		region           string
		wantRateFound    bool
	}{
		{
			name:             "Valid t3.micro",
			instanceClass:    "db.t3.micro",
			engine:           "mysql",
			multiAZ:          false,
			allocatedStorage: 20,
			storageType:      "gp2",
			region:           "us-east-1",
			wantRateFound:    true,
		},
		{
			name:             "Valid m5.large MultiAZ",
			instanceClass:    "db.m5.large",
			engine:           "postgres",
			multiAZ:          true,
			allocatedStorage: 100,
			storageType:      "io1",
			region:           "us-west-2",
			wantRateFound:    true,
		},
		{
			name:          "Invalid Instance Class",
			instanceClass: "db.invalid",
			wantRateFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetRDSInstancePricing(tt.instanceClass, tt.engine, tt.multiAZ, tt.allocatedStorage, tt.storageType, tt.region)
			if tt.wantRateFound {
				assert.Greater(t, pricing.Components[0].Rate, 0.0)
				// Storage is now handled via hidden dependency
				assert.Len(t, pricing.Components, 1)
			} else {
				assert.Equal(t, 0.0, pricing.Components[0].Rate)
			}
		})
	}
}

func TestCalculateRDSInstanceCost(t *testing.T) {
	// Base rates (approximate check)
	// db.t3.micro = 0.017
	// gp2 = 0.115/GB-month

	// Case 1: 1 hour, t3.micro, 20GB gp2
	// Instance: 0.017 * 1 = 0.017
	// Storage: 0.115 * 20 * (1/720) = 0.115 * 0.0277 = 0.00319
	// Total: 0.02019

	// Case 1: 1 hour, t3.micro
	cost := CalculateRDSInstanceCost(time.Hour, "db.t3.micro", "mysql", false, 20, "gp2", "us-east-1")
	assert.Greater(t, cost, 0.01)
	assert.Less(t, cost, 0.02) // Strictly 0.017

	// Case 2: MultiAZ -> 2x instance
	costMAZ := CalculateRDSInstanceCost(time.Hour, "db.t3.micro", "mysql", true, 20, "gp2", "us-east-1")
	assert.Greater(t, costMAZ, cost*1.8) // Should be roughly double
}
