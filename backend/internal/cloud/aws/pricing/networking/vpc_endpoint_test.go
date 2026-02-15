package networking

import (
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateVPCEndpointCost(t *testing.T) {
	tests := []struct {
		name            string
		duration        time.Duration
		endpointType    string
		dataProcessedGB float64
		numENIs         int
		region          string
		expectedCost    float64
	}{
		{
			name:            "Gateway Endpoint (Free)",
			duration:        time.Hour,
			endpointType:    "Gateway",
			dataProcessedGB: 1000,
			numENIs:         0, // Not applicable
			region:          "us-east-1",
			expectedCost:    0.0,
		},
		{
			name:            "Interface Endpoint (1 hour, 1 ENI, 0 GB)",
			duration:        time.Hour,
			endpointType:    "Interface",
			dataProcessedGB: 0,
			numENIs:         1,
			region:          "us-east-1",
			expectedCost:    0.01, // $0.01/hour
		},
		{
			name:            "Interface Endpoint (1 hour, 2 ENIs, 0 GB)",
			duration:        time.Hour,
			endpointType:    "Interface",
			dataProcessedGB: 0,
			numENIs:         2,
			region:          "us-east-1",
			expectedCost:    0.02, // 2 * $0.01/hour
		},
		{
			name:            "Interface Endpoint (1 hour, 1 ENI, 10 GB)",
			duration:        time.Hour,
			endpointType:    "Interface",
			dataProcessedGB: 10,
			numENIs:         1,
			region:          "us-east-1",
			expectedCost:    0.11, // $0.01 + (10 * $0.01)
		},
		{
			name:            "Interface Endpoint (720 hours, 1 ENI, 100 GB)",
			duration:        720 * time.Hour,
			endpointType:    "Interface",
			dataProcessedGB: 100,
			numENIs:         1,
			region:          "us-east-1",
			expectedCost:    8.20, // (720 * 0.01) + (100 * 0.01) = 7.20 + 1.00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateVPCEndpointCost(tt.duration, tt.endpointType, tt.dataProcessedGB, tt.numENIs, tt.region)
			assert.InDelta(t, tt.expectedCost, cost, 0.0001)
		})
	}
}

func TestGetVPCEndpointPricing(t *testing.T) {
	// Test Gateway
	gatewayPricing := GetVPCEndpointPricing("Gateway", "us-east-1")
	assert.Equal(t, "vpc_endpoint", gatewayPricing.ResourceType)
	assert.Len(t, gatewayPricing.Components, 1)
	assert.Equal(t, domainpricing.PerHour, gatewayPricing.Components[0].Model)
	assert.Equal(t, 0.0, gatewayPricing.Components[0].Rate)

	// Test Interface
	interfacePricing := GetVPCEndpointPricing("Interface", "us-east-1")
	assert.Equal(t, "vpc_endpoint", interfacePricing.ResourceType)
	assert.Len(t, interfacePricing.Components, 2)
	assert.Equal(t, "Interface Endpoint Hourly (per ENI)", interfacePricing.Components[0].Name)
	assert.Equal(t, 0.01, interfacePricing.Components[0].Rate)
	assert.Equal(t, "Interface Endpoint Data Processing", interfacePricing.Components[1].Name)
	assert.Equal(t, 0.01, interfacePricing.Components[1].Rate)
}
