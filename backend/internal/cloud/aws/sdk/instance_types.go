package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awsmodel "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
)

// InstanceTypeService provides instance type information with caching and fallback
type InstanceTypeService struct {
	client     *AWSClient
	cache      map[string]*awsmodel.InstanceTypeInfo
	cacheMutex sync.RWMutex
	staticData map[string]*awsmodel.InstanceTypeInfo
}

// NewInstanceTypeService creates a new instance type service
func NewInstanceTypeService(client *AWSClient) (*InstanceTypeService, error) {
	service := &InstanceTypeService{
		client:     client,
		cache:      make(map[string]*awsmodel.InstanceTypeInfo),
		staticData: make(map[string]*awsmodel.InstanceTypeInfo),
	}

	// Load static data
	if err := service.loadFromStatic(); err != nil {
		// Log error but don't fail - we can still use AWS API
		fmt.Printf("Warning: Failed to load static instance type data: %v\n", err)
	}

	return service, nil
}

// GetInstanceType retrieves information about a specific instance type
// Checks cache first, then AWS API, then static data
func (s *InstanceTypeService) GetInstanceType(ctx context.Context, name string, region string) (*awsmodel.InstanceTypeInfo, error) {
	// Normalize instance type name
	name = normalizeInstanceTypeName(name)

	// Check cache first
	s.cacheMutex.RLock()
	if cached, ok := s.cache[name]; ok {
		s.cacheMutex.RUnlock()
		if cached.IsAvailableInRegion(region) {
			return cached, nil
		}
	} else {
		s.cacheMutex.RUnlock()
	}

	// Try AWS API if client is available
	if s.client != nil && s.client.EC2 != nil {
		info, err := s.loadFromAWS(ctx, name, region)
		if err == nil && info != nil {
			// Cache the result
			s.cacheMutex.Lock()
			s.cache[name] = info
			s.cacheMutex.Unlock()
			return info, nil
		}
	}

	// Fall back to static data
	s.cacheMutex.RLock()
	if static, ok := s.staticData[name]; ok {
		s.cacheMutex.RUnlock()
		// Cache it for future use
		s.cacheMutex.Lock()
		s.cache[name] = static
		s.cacheMutex.Unlock()
		return static, nil
	}
	s.cacheMutex.RUnlock()

	return nil, awsmodel.ErrInstanceTypeNotFound
}

// ListInstanceTypes lists all available instance types
func (s *InstanceTypeService) ListInstanceTypes(ctx context.Context, region string, filters *awsmodel.InstanceTypeFilters) ([]*awsmodel.InstanceTypeInfo, error) {
	var results []*awsmodel.InstanceTypeInfo

	// Try to load from AWS API first
	if s.client != nil && s.client.EC2 != nil {
		awsTypes, err := s.loadAllFromAWS(ctx, region)
		if err == nil {
			results = awsTypes
		}
	}

	// If AWS API didn't return results, use static data
	if len(results) == 0 {
		s.cacheMutex.RLock()
		for _, info := range s.staticData {
			results = append(results, info)
		}
		s.cacheMutex.RUnlock()
	}

	// Apply filters if provided
	if filters != nil {
		results = s.applyFilters(results, filters)
	}

	return results, nil
}

// ListByCategory returns all instance types in a specific category
func (s *InstanceTypeService) ListByCategory(ctx context.Context, category awsmodel.InstanceCategory, region string) ([]*awsmodel.InstanceTypeInfo, error) {
	allTypes, err := s.ListInstanceTypes(ctx, region, nil)
	if err != nil {
		return nil, err
	}

	var filtered []*awsmodel.InstanceTypeInfo
	for _, info := range allTypes {
		if info.Category == category && info.IsAvailableInRegion(region) {
			filtered = append(filtered, info)
		}
	}

	return filtered, nil
}

// ListFreeTier returns all free tier eligible instance types
func (s *InstanceTypeService) ListFreeTier(ctx context.Context, region string) ([]*awsmodel.InstanceTypeInfo, error) {
	allTypes, err := s.ListInstanceTypes(ctx, region, nil)
	if err != nil {
		return nil, err
	}

	var freeTier []*awsmodel.InstanceTypeInfo
	for _, info := range allTypes {
		if info.FreeTierEligible && info.IsAvailableInRegion(region) {
			freeTier = append(freeTier, info)
		}
	}

	return freeTier, nil
}

// RefreshCache forces a refresh of the cache from AWS API
func (s *InstanceTypeService) RefreshCache(ctx context.Context, region string) error {
	if s.client == nil || s.client.EC2 == nil {
		return fmt.Errorf("AWS client not available")
	}

	types, err := s.loadAllFromAWS(ctx, region)
	if err != nil {
		return err
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for _, info := range types {
		s.cache[normalizeInstanceTypeName(info.Name)] = info
	}

	return nil
}

// loadFromAWS fetches a specific instance type from AWS API
func (s *InstanceTypeService) loadFromAWS(ctx context.Context, name string, region string) (*awsmodel.InstanceTypeInfo, error) {
	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []types.InstanceType{
			types.InstanceType(name),
		},
	}

	output, err := s.client.EC2.DescribeInstanceTypes(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(output.InstanceTypes) == 0 {
		return nil, awsmodel.ErrInstanceTypeNotFound
	}

	return s.convertAWSInstanceType(output.InstanceTypes[0], region), nil
}

// loadAllFromAWS fetches all instance types from AWS API
func (s *InstanceTypeService) loadAllFromAWS(ctx context.Context, region string) ([]*awsmodel.InstanceTypeInfo, error) {
	var allTypes []*awsmodel.InstanceTypeInfo
	var nextToken *string

	for {
		input := &ec2.DescribeInstanceTypesInput{
			MaxResults: aws.Int32(100),
		}
		if nextToken != nil {
			input.NextToken = nextToken
		}

		output, err := s.client.EC2.DescribeInstanceTypes(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, awsType := range output.InstanceTypes {
			allTypes = append(allTypes, s.convertAWSInstanceType(awsType, region))
		}

		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return allTypes, nil
}

// convertAWSInstanceType converts AWS SDK instance type to our model
func (s *InstanceTypeService) convertAWSInstanceType(awsType types.InstanceTypeInfo, region string) *awsmodel.InstanceTypeInfo {
	name := string(awsType.InstanceType)
	category := awsmodel.CategorizeInstanceType(name)

	info := &awsmodel.InstanceTypeInfo{
		Name:                        name,
		Category:                    category,
		FreeTierEligible:            awsmodel.IsFreeTierEligible(name, region),
		SupportedArchitectures:      []string{},
		SupportedVirtualizationTypes: []string{},
		Region:                      region,
	}

	// VCPU
	if awsType.VCpuInfo != nil {
		info.VCPU = int(aws.ToInt32(awsType.VCpuInfo.DefaultVCpus))
	}

	// Memory
	if awsType.MemoryInfo != nil && awsType.MemoryInfo.SizeInMiB != nil {
		info.MemoryGiB = float64(*awsType.MemoryInfo.SizeInMiB) / 1024.0
	}

	// Storage
	if awsType.InstanceStorageInfo != nil {
		info.HasLocalStorage = awsType.InstanceStorageInfo.TotalSizeInGB != nil && *awsType.InstanceStorageInfo.TotalSizeInGB > 0
		if info.HasLocalStorage && awsType.InstanceStorageInfo.TotalSizeInGB != nil {
			size := float64(*awsType.InstanceStorageInfo.TotalSizeInGB)
			info.LocalStorageSizeGiB = &size
		}
		// Determine storage type from disks if available
		if len(awsType.InstanceStorageInfo.Disks) > 0 {
			// Use the first disk's type as the storage type
			info.StorageType = "NVMe SSD" // Most common for local storage
		} else if info.HasLocalStorage {
			info.StorageType = "Instance Store"
		}
	} else {
		info.StorageType = "EBS"
		info.HasLocalStorage = false
	}

	// Network
	if awsType.NetworkInfo != nil {
		if awsType.NetworkInfo.NetworkPerformance != nil {
			// Parse network performance (e.g., "Up to 5 Gigabit", "10 Gigabit")
			info.MaxNetworkGbps = parseNetworkPerformance(*awsType.NetworkInfo.NetworkPerformance)
		}
	}

	// EBS Bandwidth
	if awsType.EbsInfo != nil {
		if awsType.EbsInfo.EbsOptimizedSupport != types.EbsOptimizedSupportUnsupported {
			if awsType.EbsInfo.EbsOptimizedInfo != nil && awsType.EbsInfo.EbsOptimizedInfo.BaselineBandwidthInMbps != nil {
				bandwidth := float64(*awsType.EbsInfo.EbsOptimizedInfo.BaselineBandwidthInMbps) / 1000.0
				info.EBSBandwidthGbps = &bandwidth
			}
		}
	}

	// Supported architectures
	if awsType.ProcessorInfo != nil && awsType.ProcessorInfo.SupportedArchitectures != nil {
		for _, arch := range awsType.ProcessorInfo.SupportedArchitectures {
			info.SupportedArchitectures = append(info.SupportedArchitectures, string(arch))
		}
	}

	// Supported virtualization types
	if awsType.SupportedVirtualizationTypes != nil {
		for _, vt := range awsType.SupportedVirtualizationTypes {
			info.SupportedVirtualizationTypes = append(info.SupportedVirtualizationTypes, string(vt))
		}
	}

	return info
}

// loadFromStatic loads instance types from the static JSON file
func (s *InstanceTypeService) loadFromStatic() error {
	// Try to find the JSON file relative to this package
	// The file is in backend/internal/cloud/aws/models/compute/data/instance_types.json
	wd, _ := os.Getwd()
	possiblePaths := []string{
		filepath.Join(wd, "internal", "cloud", "aws", "models", "compute", "data", "instance_types.json"),
		filepath.Join(wd, "..", "..", "..", "..", "..", "models", "compute", "data", "instance_types.json"),
		filepath.Join(wd, "backend", "internal", "cloud", "aws", "models", "compute", "data", "instance_types.json"),
	}

	var jsonData []byte
	var err error
	for _, path := range possiblePaths {
		jsonData, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to find instance_types.json: %w", err)
	}

	var types []*awsmodel.InstanceTypeInfo
	if err := json.Unmarshal(jsonData, &types); err != nil {
		return fmt.Errorf("failed to parse instance_types.json: %w", err)
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for _, info := range types {
		normalizedName := normalizeInstanceTypeName(info.Name)
		s.staticData[normalizedName] = info
	}

	return nil
}

// applyFilters applies filters to the instance type list
func (s *InstanceTypeService) applyFilters(types []*awsmodel.InstanceTypeInfo, filters *awsmodel.InstanceTypeFilters) []*awsmodel.InstanceTypeInfo {
	var filtered []*awsmodel.InstanceTypeInfo

	for _, info := range types {
		// Category filter
		if filters.Category != nil && info.Category != *filters.Category {
			continue
		}

		// VCPU filters
		if filters.MinVCPU != nil && info.VCPU < *filters.MinVCPU {
			continue
		}
		if filters.MaxVCPU != nil && info.VCPU > *filters.MaxVCPU {
			continue
		}

		// Memory filters
		if filters.MinMemoryGiB != nil && info.MemoryGiB < *filters.MinMemoryGiB {
			continue
		}
		if filters.MaxMemoryGiB != nil && info.MemoryGiB > *filters.MaxMemoryGiB {
			continue
		}

		// Free tier filter
		if filters.FreeTierOnly && !info.FreeTierEligible {
			continue
		}

		// Local storage filter
		if filters.HasLocalStorage != nil && info.HasLocalStorage != *filters.HasLocalStorage {
			continue
		}

		// Architecture filter
		if filters.SupportedArchitecture != nil {
			found := false
			for _, arch := range info.SupportedArchitectures {
				if arch == *filters.SupportedArchitecture {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Region filter
		if filters.Region != "" && !info.IsAvailableInRegion(filters.Region) {
			continue
		}

		filtered = append(filtered, info)
	}

	return filtered
}

// normalizeInstanceTypeName normalizes instance type name to lowercase
func normalizeInstanceTypeName(name string) string {
	return strings.ToLower(name)
}

// parseNetworkPerformance parses AWS network performance string to Gbps
func parseNetworkPerformance(perf string) float64 {
	// AWS formats: "Up to 5 Gigabit", "10 Gigabit", "25 Gigabit", etc.
	// This is a simplified parser - in production, you might want more robust parsing
	perf = strings.ToLower(perf)
	
	// Extract numbers and convert to Gbps
	// For now, return a default value - this would need more sophisticated parsing
	if strings.Contains(perf, "25") {
		return 25.0
	}
	if strings.Contains(perf, "10") {
		return 10.0
	}
	if strings.Contains(perf, "5") {
		return 5.0
	}
	if strings.Contains(perf, "1") {
		return 1.0
	}
	return 0.0
}
