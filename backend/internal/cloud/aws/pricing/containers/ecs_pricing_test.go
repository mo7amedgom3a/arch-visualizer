package containers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFargateCost(t *testing.T) {
	tests := []struct {
		name     string
		vcpu     float64
		memoryGB float64
		duration time.Duration
		region   string
		spot     bool
		wantMin  float64
		wantMax  float64
	}{
		{
			name:     "1 vCPU, 2GB memory, 1 hour, us-east-1",
			vcpu:     1.0,
			memoryGB: 2.0,
			duration: time.Hour,
			region:   "us-east-1",
			spot:     false,
			wantMin:  0.04,
			wantMax:  0.06,
		},
		{
			name:     "0.25 vCPU, 0.5GB memory, 1 hour",
			vcpu:     0.25,
			memoryGB: 0.5,
			duration: time.Hour,
			region:   "us-east-1",
			spot:     false,
			wantMin:  0.01,
			wantMax:  0.02,
		},
		{
			name:     "4 vCPU, 8GB memory, 24 hours",
			vcpu:     4.0,
			memoryGB: 8.0,
			duration: 24 * time.Hour,
			region:   "us-east-1",
			spot:     false,
			wantMin:  4.0,
			wantMax:  5.0,
		},
		{
			name:     "Spot instance should be cheaper",
			vcpu:     1.0,
			memoryGB: 2.0,
			duration: time.Hour,
			region:   "us-east-1",
			spot:     true,
			wantMin:  0.01,
			wantMax:  0.02,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateFargateCost(tt.vcpu, tt.memoryGB, tt.duration, tt.region, tt.spot)
			assert.GreaterOrEqual(t, got, tt.wantMin)
			assert.LessOrEqual(t, got, tt.wantMax)
		})
	}
}

func TestGetFargatePricing(t *testing.T) {
	tests := []struct {
		name     string
		vcpu     float64
		memoryGB float64
		region   string
		spot     bool
	}{
		{
			name:     "On-demand pricing",
			vcpu:     1.0,
			memoryGB: 2.0,
			region:   "us-east-1",
			spot:     false,
		},
		{
			name:     "Spot pricing",
			vcpu:     2.0,
			memoryGB: 4.0,
			region:   "us-west-2",
			spot:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetFargatePricing(tt.vcpu, tt.memoryGB, tt.region, tt.spot)

			assert.NotNil(t, pricing)
			assert.Equal(t, "ecs_fargate", pricing.ResourceType)
			assert.Len(t, pricing.Components, 2) // vCPU and Memory components

			// Verify vCPU component
			assert.Equal(t, "Fargate vCPU", pricing.Components[0].Name)
			assert.Greater(t, pricing.Components[0].Rate, 0.0)

			// Verify Memory component
			assert.Equal(t, "Fargate Memory", pricing.Components[1].Name)
			assert.Greater(t, pricing.Components[1].Rate, 0.0)

			// Check metadata
			assert.NotNil(t, pricing.Metadata)
			assert.Equal(t, tt.vcpu, pricing.Metadata["vcpu"])
			assert.Equal(t, tt.memoryGB, pricing.Metadata["memory_gb"])
		})
	}
}

func TestGetECSTaskDefinitionPricing(t *testing.T) {
	tests := []struct {
		name       string
		cpu        string
		memory     string
		region     string
		launchType string
		wantType   string
	}{
		{
			name:       "Fargate task definition",
			cpu:        "256",
			memory:     "512",
			region:     "us-east-1",
			launchType: "FARGATE",
			wantType:   "ecs_fargate",
		},
		{
			name:       "EC2 task definition",
			cpu:        "1024",
			memory:     "2048",
			region:     "us-east-1",
			launchType: "EC2",
			wantType:   "ecs_task_definition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetECSTaskDefinitionPricing(tt.cpu, tt.memory, tt.region, tt.launchType)

			assert.NotNil(t, pricing)
			assert.Equal(t, tt.wantType, pricing.ResourceType)
		})
	}
}

func TestParseCPUToVCPU(t *testing.T) {
	tests := []struct {
		cpu      string
		expected float64
	}{
		{"256", 0.25},
		{"512", 0.5},
		{"1024", 1.0},
		{"2048", 2.0},
		{"4096", 4.0},
		{"unknown", 0.25}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.cpu, func(t *testing.T) {
			got := parseCPUToVCPU(tt.cpu)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetECSServicePricing(t *testing.T) {
	t.Run("Fargate service with multiple tasks", func(t *testing.T) {
		pricing := GetECSServicePricing("my-service", "FARGATE", "256", "512", "us-east-1", 3)

		assert.NotNil(t, pricing)
		assert.Equal(t, "ecs_service", pricing.ResourceType)
		assert.Equal(t, "my-service", pricing.Metadata["service_name"])
		assert.Equal(t, 3, pricing.Metadata["desired_count"])
	})

	t.Run("EC2 service", func(t *testing.T) {
		pricing := GetECSServicePricing("my-service", "EC2", "1024", "2048", "us-east-1", 2)

		assert.NotNil(t, pricing)
		assert.Equal(t, "ecs_service", pricing.ResourceType)
		assert.Equal(t, "EC2", pricing.Metadata["launch_type"])
	})
}

func TestGetECSClusterPricing(t *testing.T) {
	pricing := GetECSClusterPricing("my-cluster", "us-east-1")

	assert.NotNil(t, pricing)
	assert.Equal(t, "ecs_cluster", pricing.ResourceType)
	assert.Equal(t, "my-cluster", pricing.Metadata["cluster_name"])
	assert.Empty(t, pricing.Components) // Clusters have no direct charge
}
