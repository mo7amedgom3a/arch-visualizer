package seed

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
)

// DatabaseQueryService provides functions to query database and convert to domain models
type DatabaseQueryService struct {
	resourceRepo *resourcerepo.ResourceRepository
	mapper       *ResourceMapper
}

// NewDatabaseQueryService creates a new database query service
func NewDatabaseQueryService() (*DatabaseQueryService, error) {
	resourceRepo, err := resourcerepo.NewResourceRepository(slog.Default())
	if err != nil {
		return nil, fmt.Errorf("failed to create resource repository: %w", err)
	}

	return &DatabaseQueryService{
		resourceRepo: resourceRepo,
		mapper:       NewResourceMapper(),
	}, nil
}

// GetProjectResourcesAsDomainModels queries all resources for a project and converts them to domain models
func (dqs *DatabaseQueryService) GetProjectResourcesAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]interface{}, error) {
	fmt.Printf("[LOG] Querying all resources for ProjectID: %s\n", projectID.String())
	// Query database
	resources, err := dqs.resourceRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		fmt.Printf("[LOG] Failed to query resources for ProjectID %s: %v\n", projectID.String(), err)
		return nil, fmt.Errorf("failed to query resources: %w", err)
	}

	fmt.Printf("[LOG] Found %d resources in database for project %s. Converting to domain models...\n", len(resources), projectID.String())

	// Convert each resource to domain model
	domainResources := make([]interface{}, 0, len(resources))
	for _, resource := range resources {
		// Preload ResourceType if not already loaded
		if resource.ResourceType.Name == "" {
			fmt.Printf("[LOG] ResourceType not loaded for resource %s (ID: %s), re-fetching with preload...\n", resource.Name, resource.ID)
			// Re-fetch with preload if needed
			loadedResource, err := dqs.resourceRepo.FindByID(ctx, resource.ID)
			if err != nil {
				fmt.Printf("  ⚠ Warning: failed to load resource %s: %v\n", resource.ID, err)
				continue
			}
			resource = loadedResource
		}

		domainResource, err := dqs.mapper.DatabaseResourceToDomain(resource)
		if err != nil {
			fmt.Printf("  ⚠ Warning: failed to convert resource %s (%s): %v\n", resource.Name, resource.ResourceType.Name, err)
			continue
		} else {
			fmt.Printf("[LOG] Mapped resource '%s' (type: %s, db id: %s) to domain model type: %T\n",
				resource.Name, resource.ResourceType.Name, resource.ID, domainResource)
		}

		domainResources = append(domainResources, domainResource)
	}

	fmt.Printf("[LOG] Finished mapping %d resources from database to domain models for project %s\n", len(domainResources), projectID.String())
	return domainResources, nil
}

// GetResourceByIDAsDomainModel queries a single resource by ID and converts it to domain model
func (dqs *DatabaseQueryService) GetResourceByIDAsDomainModel(ctx context.Context, resourceID uuid.UUID) (interface{}, error) {
	fmt.Printf("[LOG] Querying resource by ID: %s\n", resourceID.String())
	// Query database
	resource, err := dqs.resourceRepo.FindByID(ctx, resourceID)
	if err != nil {
		fmt.Printf("[LOG] Failed to query resource by ID %s: %v\n", resourceID.String(), err)
		return nil, fmt.Errorf("failed to query resource: %w", err)
	}

	fmt.Printf("[LOG] Converting queried resource '%s' (type: %s) to domain model\n", resource.Name, resource.ResourceType.Name)
	// Convert to domain model
	domainResource, err := dqs.mapper.DatabaseResourceToDomain(resource)
	if err != nil {
		fmt.Printf("[LOG] Failed to convert database resource '%s' to domain model: %v\n", resource.Name, err)
		return nil, fmt.Errorf("failed to convert resource to domain model: %w", err)
	}
	fmt.Printf("[LOG] Mapped resource '%s' (db id: %s) to domain model type: %T\n", resource.Name, resource.ID, domainResource)
	return domainResource, nil
}

// GetResourcesByTypeAsDomainModels queries resources by type and converts them to domain models
func (dqs *DatabaseQueryService) GetResourcesByTypeAsDomainModels(ctx context.Context, projectID uuid.UUID, resourceTypeID uint) ([]interface{}, error) {
	fmt.Printf("[LOG] Querying resources for project %s by type id %d\n", projectID.String(), resourceTypeID)
	// Query database
	resources, err := dqs.resourceRepo.FindByProjectIDAndType(ctx, projectID, resourceTypeID)
	if err != nil {
		fmt.Printf("[LOG] Failed to query resources by type for project %s: %v\n", projectID.String(), err)
		return nil, fmt.Errorf("failed to query resources: %w", err)
	}

	fmt.Printf("[LOG] Found %d resources in DB for project %s by type id %d. Mapping to domain models...\n", len(resources), projectID.String(), resourceTypeID)

	// Convert each resource to domain model
	domainResources := make([]interface{}, 0, len(resources))
	for _, resource := range resources {
		domainResource, err := dqs.mapper.DatabaseResourceToDomain(resource)
		if err != nil {
			fmt.Printf("  ⚠ Warning: failed to convert resource %s: %v\n", resource.Name, err)
			continue
		}
		fmt.Printf("[LOG] Mapped resource '%s' (type: %s, db id: %s) to domain model type: %T\n",
			resource.Name, resource.ResourceType.Name, resource.ID, domainResource)
		domainResources = append(domainResources, domainResource)
	}

	fmt.Printf("[LOG] Finished mapping %d resources (by type) to domain models for project %s\n", len(domainResources), projectID.String())
	return domainResources, nil
}

// Typed query functions for specific resource types

// GetVPCsAsDomainModels returns all VPCs from a project as domain models
func (dqs *DatabaseQueryService) GetVPCsAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domainnetworking.VPC, error) {
	fmt.Printf("[LOG] Getting all VPCs as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	vpcs := make([]*domainnetworking.VPC, 0)
	for _, resource := range domainResources {
		if vpc, ok := resource.(*domainnetworking.VPC); ok {
			vpcs = append(vpcs, vpc)
		}
	}
	fmt.Printf("[LOG] %d VPC domain models found for project %s\n", len(vpcs), projectID.String())
	return vpcs, nil
}

// GetSubnetsAsDomainModels returns all Subnets from a project as domain models
func (dqs *DatabaseQueryService) GetSubnetsAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domainnetworking.Subnet, error) {
	fmt.Printf("[LOG] Getting all Subnets as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	subnets := make([]*domainnetworking.Subnet, 0)
	for _, resource := range domainResources {
		if subnet, ok := resource.(*domainnetworking.Subnet); ok {
			subnets = append(subnets, subnet)
		}
	}
	fmt.Printf("[LOG] %d Subnet domain models found for project %s\n", len(subnets), projectID.String())
	return subnets, nil
}

// GetInstancesAsDomainModels returns all EC2 instances from a project as domain models
func (dqs *DatabaseQueryService) GetInstancesAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domaincompute.Instance, error) {
	fmt.Printf("[LOG] Getting all EC2 Instances as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	instances := make([]*domaincompute.Instance, 0)
	for _, resource := range domainResources {
		if instance, ok := resource.(*domaincompute.Instance); ok {
			instances = append(instances, instance)
		}
	}
	fmt.Printf("[LOG] %d EC2 Instance domain models found for project %s\n", len(instances), projectID.String())
	return instances, nil
}

// GetSecurityGroupsAsDomainModels returns all Security Groups from a project as domain models
func (dqs *DatabaseQueryService) GetSecurityGroupsAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domainnetworking.SecurityGroup, error) {
	fmt.Printf("[LOG] Getting all Security Groups as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	securityGroups := make([]*domainnetworking.SecurityGroup, 0)
	for _, resource := range domainResources {
		if sg, ok := resource.(*domainnetworking.SecurityGroup); ok {
			securityGroups = append(securityGroups, sg)
		}
	}
	fmt.Printf("[LOG] %d Security Group domain models found for project %s\n", len(securityGroups), projectID.String())
	return securityGroups, nil
}

// GetLoadBalancersAsDomainModels returns all Load Balancers from a project as domain models
func (dqs *DatabaseQueryService) GetLoadBalancersAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domaincompute.LoadBalancer, error) {
	fmt.Printf("[LOG] Getting all Load Balancers as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	loadBalancers := make([]*domaincompute.LoadBalancer, 0)
	for _, resource := range domainResources {
		if lb, ok := resource.(*domaincompute.LoadBalancer); ok {
			loadBalancers = append(loadBalancers, lb)
		}
	}
	fmt.Printf("[LOG] %d Load Balancer domain models found for project %s\n", len(loadBalancers), projectID.String())
	return loadBalancers, nil
}

// GetAutoScalingGroupsAsDomainModels returns all Auto Scaling Groups from a project as domain models
func (dqs *DatabaseQueryService) GetAutoScalingGroupsAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domaincompute.AutoScalingGroup, error) {
	fmt.Printf("[LOG] Getting all Auto Scaling Groups as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	asgs := make([]*domaincompute.AutoScalingGroup, 0)
	for _, resource := range domainResources {
		if asg, ok := resource.(*domaincompute.AutoScalingGroup); ok {
			asgs = append(asgs, asg)
		}
	}
	fmt.Printf("[LOG] %d Auto Scaling Group domain models found for project %s\n", len(asgs), projectID.String())
	return asgs, nil
}

// GetLambdaFunctionsAsDomainModels returns all Lambda functions from a project as domain models
func (dqs *DatabaseQueryService) GetLambdaFunctionsAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domaincompute.LambdaFunction, error) {
	fmt.Printf("[LOG] Getting all Lambda Functions as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	lambdas := make([]*domaincompute.LambdaFunction, 0)
	for _, resource := range domainResources {
		if lambda, ok := resource.(*domaincompute.LambdaFunction); ok {
			lambdas = append(lambdas, lambda)
		}
	}
	fmt.Printf("[LOG] %d Lambda Function domain models found for project %s\n", len(lambdas), projectID.String())
	return lambdas, nil
}

// GetS3BucketsAsDomainModels returns all S3 buckets from a project as domain models
func (dqs *DatabaseQueryService) GetS3BucketsAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domainstorage.S3Bucket, error) {
	fmt.Printf("[LOG] Getting all S3 Buckets as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	buckets := make([]*domainstorage.S3Bucket, 0)
	for _, resource := range domainResources {
		if bucket, ok := resource.(*domainstorage.S3Bucket); ok {
			buckets = append(buckets, bucket)
		}
	}
	fmt.Printf("[LOG] %d S3 Bucket domain models found for project %s\n", len(buckets), projectID.String())
	return buckets, nil
}

// GetEBSVolumesAsDomainModels returns all EBS volumes from a project as domain models
func (dqs *DatabaseQueryService) GetEBSVolumesAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domainstorage.EBSVolume, error) {
	fmt.Printf("[LOG] Getting all EBS Volumes as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	volumes := make([]*domainstorage.EBSVolume, 0)
	for _, resource := range domainResources {
		if volume, ok := resource.(*domainstorage.EBSVolume); ok {
			volumes = append(volumes, volume)
		}
	}
	fmt.Printf("[LOG] %d EBS Volume domain models found for project %s\n", len(volumes), projectID.String())
	return volumes, nil
}

// GetNATGatewaysAsDomainModels returns all NAT Gateways from a project as domain models
func (dqs *DatabaseQueryService) GetNATGatewaysAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domainnetworking.NATGateway, error) {
	fmt.Printf("[LOG] Getting all NAT Gateways as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	natGateways := make([]*domainnetworking.NATGateway, 0)
	for _, resource := range domainResources {
		if nat, ok := resource.(*domainnetworking.NATGateway); ok {
			natGateways = append(natGateways, nat)
		}
	}
	fmt.Printf("[LOG] %d NAT Gateway domain models found for project %s\n", len(natGateways), projectID.String())
	return natGateways, nil
}

// GetInternetGatewaysAsDomainModels returns all Internet Gateways from a project as domain models
func (dqs *DatabaseQueryService) GetInternetGatewaysAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domainnetworking.InternetGateway, error) {
	fmt.Printf("[LOG] Getting all Internet Gateways as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	igws := make([]*domainnetworking.InternetGateway, 0)
	for _, resource := range domainResources {
		if igw, ok := resource.(*domainnetworking.InternetGateway); ok {
			igws = append(igws, igw)
		}
	}
	fmt.Printf("[LOG] %d Internet Gateway domain models found for project %s\n", len(igws), projectID.String())
	return igws, nil
}

// GetElasticIPsAsDomainModels returns all Elastic IPs from a project as domain models
func (dqs *DatabaseQueryService) GetElasticIPsAsDomainModels(ctx context.Context, projectID uuid.UUID) ([]*domainnetworking.ElasticIP, error) {
	fmt.Printf("[LOG] Getting all Elastic IPs as domain models for project %s\n", projectID.String())
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		return nil, err
	}

	eips := make([]*domainnetworking.ElasticIP, 0)
	for _, resource := range domainResources {
		if eip, ok := resource.(*domainnetworking.ElasticIP); ok {
			eips = append(eips, eip)
		}
	}
	fmt.Printf("[LOG] %d Elastic IP domain models found for project %s\n", len(eips), projectID.String())
	return eips, nil
}

// PrintProjectResourcesSummary prints a summary of all resources in a project as domain models
func (dqs *DatabaseQueryService) PrintProjectResourcesSummary(ctx context.Context, projectID uuid.UUID) error {
	fmt.Printf("[LOG] Printing project resources summary for project %s\n", projectID.String())
	fmt.Println("\n📊 Project Resources Summary (as Domain Models):")

	// Get all resources
	domainResources, err := dqs.GetProjectResourcesAsDomainModels(ctx, projectID)
	if err != nil {
		fmt.Printf("[LOG] Failed to get domain resources for summary for project %s: %v\n", projectID.String(), err)
		return fmt.Errorf("failed to get resources: %w", err)
	}
	fmt.Printf("[LOG] Total domain models counted for summary: %d for project %s\n", len(domainResources), projectID.String())

	// Count by type
	typeCounts := make(map[string]int)
	for _, resource := range domainResources {
		typeName := fmt.Sprintf("%T", resource)
		typeCounts[typeName]++
	}

	fmt.Printf("  Total Resources: %d\n", len(domainResources))
	fmt.Println("\n  Resource Types:")
	for typeName, count := range typeCounts {
		fmt.Printf("    - %s: %d\n", typeName, count)
	}

	// Print details for each type
	fmt.Println("\n  Resource Details:")

	// VPCs
	vpcs, _ := dqs.GetVPCsAsDomainModels(ctx, projectID)
	if len(vpcs) > 0 {
		fmt.Printf("\n  VPCs (%d):\n", len(vpcs))
		for _, vpc := range vpcs {
			fmt.Printf("    - %s (ID: %s, CIDR: %s)\n", vpc.Name, vpc.ID, vpc.CIDR)
		}
	}

	// Subnets
	subnets, _ := dqs.GetSubnetsAsDomainModels(ctx, projectID)
	if len(subnets) > 0 {
		fmt.Printf("\n  Subnets (%d):\n", len(subnets))
		for _, subnet := range subnets {
			fmt.Printf("    - %s (ID: %s, CIDR: %s)\n", subnet.Name, subnet.ID, subnet.CIDR)
		}
	}

	// Instances
	instances, _ := dqs.GetInstancesAsDomainModels(ctx, projectID)
	if len(instances) > 0 {
		fmt.Printf("\n  EC2 Instances (%d):\n", len(instances))
		for _, instance := range instances {
			fmt.Printf("    - %s (ID: %s, Type: %s)\n", instance.Name, instance.ID, instance.InstanceType)
		}
	}

	// Security Groups
	sgs, _ := dqs.GetSecurityGroupsAsDomainModels(ctx, projectID)
	if len(sgs) > 0 {
		fmt.Printf("\n  Security Groups (%d):\n", len(sgs))
		for _, sg := range sgs {
			fmt.Printf("    - %s (ID: %s)\n", sg.Name, sg.ID)
		}
	}

	// Load Balancers
	lbs, _ := dqs.GetLoadBalancersAsDomainModels(ctx, projectID)
	if len(lbs) > 0 {
		fmt.Printf("\n  Load Balancers (%d):\n", len(lbs))
		for _, lb := range lbs {
			fmt.Printf("    - %s (ID: %s, Type: %s)\n", lb.Name, lb.ID, lb.Type)
		}
	}

	// Auto Scaling Groups
	asgs, _ := dqs.GetAutoScalingGroupsAsDomainModels(ctx, projectID)
	if len(asgs) > 0 {
		fmt.Printf("\n  Auto Scaling Groups (%d):\n", len(asgs))
		for _, asg := range asgs {
			fmt.Printf("    - %s (ID: %s, Min: %d, Max: %d)\n", asg.Name, asg.ID, asg.MinSize, asg.MaxSize)
		}
	}

	// Lambda Functions
	lambdas, _ := dqs.GetLambdaFunctionsAsDomainModels(ctx, projectID)
	if len(lambdas) > 0 {
		fmt.Printf("\n  Lambda Functions (%d):\n", len(lambdas))
		for _, lambda := range lambdas {
			fmt.Printf("    - %s (ARN: %s)\n", lambda.FunctionName, getARN(lambda.ARN))
		}
	}

	// S3 Buckets
	buckets, _ := dqs.GetS3BucketsAsDomainModels(ctx, projectID)
	if len(buckets) > 0 {
		fmt.Printf("\n  S3 Buckets (%d):\n", len(buckets))
		for _, bucket := range buckets {
			fmt.Printf("    - %s (ID: %s, Region: %s)\n", bucket.Name, bucket.ID, bucket.Region)
		}
	}

	// NAT Gateways
	nats, _ := dqs.GetNATGatewaysAsDomainModels(ctx, projectID)
	if len(nats) > 0 {
		fmt.Printf("\n  NAT Gateways (%d):\n", len(nats))
		for _, nat := range nats {
			fmt.Printf("    - %s (ID: %s)\n", nat.Name, nat.ID)
		}
	}

	// Internet Gateways
	igws, _ := dqs.GetInternetGatewaysAsDomainModels(ctx, projectID)
	if len(igws) > 0 {
		fmt.Printf("\n  Internet Gateways (%d):\n", len(igws))
		for _, igw := range igws {
			fmt.Printf("    - %s (ID: %s)\n", igw.Name, igw.ID)
		}
	}

	fmt.Printf("[LOG] Project summary printed for project %s\n", projectID.String())
	return nil
}
