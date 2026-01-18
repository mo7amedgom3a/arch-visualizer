package compute

import (
	"context"
	"fmt"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/compute"
	awsInstanceTypes "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// AWSComputeAdapter adapts AWS-specific compute service to domain compute service
// This implements the Adapter pattern, allowing the domain layer to work with cloud-specific implementations
type AWSComputeAdapter struct {
	awsService awsservice.AWSComputeService
}

// NewAWSComputeAdapter creates a new AWS compute adapter
func NewAWSComputeAdapter(awsService awsservice.AWSComputeService) domaincompute.ComputeService {
	return &AWSComputeAdapter{
		awsService: awsService,
	}
}

// Ensure AWSComputeAdapter implements ComputeService
var _ domaincompute.ComputeService = (*AWSComputeAdapter)(nil)

// Instance Operations

func (a *AWSComputeAdapter) CreateInstance(ctx context.Context, instance *domaincompute.Instance) (*domaincompute.Instance, error) {
	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsInstance := awsmapper.FromDomainInstance(instance)
	if err := awsInstance.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsInstanceOutput, err := a.awsService.CreateInstance(ctx, awsInstance)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceFromOutput(awsInstanceOutput), nil
}

func (a *AWSComputeAdapter) GetInstance(ctx context.Context, id string) (*domaincompute.Instance, error) {
	awsInstanceOutput, err := a.awsService.GetInstance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceFromOutput(awsInstanceOutput), nil
}

func (a *AWSComputeAdapter) UpdateInstance(ctx context.Context, instance *domaincompute.Instance) (*domaincompute.Instance, error) {
	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsInstance := awsmapper.FromDomainInstance(instance)
	if err := awsInstance.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsInstanceOutput, err := a.awsService.UpdateInstance(ctx, awsInstance)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceFromOutput(awsInstanceOutput), nil
}

func (a *AWSComputeAdapter) DeleteInstance(ctx context.Context, id string) error {
	if err := a.awsService.DeleteInstance(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListInstances(ctx context.Context, filters map[string]string) ([]*domaincompute.Instance, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsInstanceOutputs, err := a.awsService.ListInstances(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainInstances := make([]*domaincompute.Instance, len(awsInstanceOutputs))
	for i, awsInstanceOutput := range awsInstanceOutputs {
		domainInstances[i] = awsmapper.ToDomainInstanceFromOutput(awsInstanceOutput)
	}

	return domainInstances, nil
}

// Instance Lifecycle Operations

func (a *AWSComputeAdapter) StartInstance(ctx context.Context, id string) error {
	if err := a.awsService.StartInstance(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) StopInstance(ctx context.Context, id string) error {
	if err := a.awsService.StopInstance(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) RebootInstance(ctx context.Context, id string) error {
	if err := a.awsService.RebootInstance(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

// Instance Type Operations

func (a *AWSComputeAdapter) GetInstanceTypeInfo(ctx context.Context, instanceType string, region string) (*awsInstanceTypes.InstanceTypeInfo, error) {
	info, err := a.awsService.GetInstanceTypeInfo(ctx, instanceType, region)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return info, nil
}

func (a *AWSComputeAdapter) ListInstanceTypesByCategory(ctx context.Context, category awsInstanceTypes.InstanceCategory, region string) ([]*awsInstanceTypes.InstanceTypeInfo, error) {
	types, err := a.awsService.ListInstanceTypesByCategory(ctx, category, region)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return types, nil
}

// Launch Template Operations

func (a *AWSComputeAdapter) CreateLaunchTemplate(ctx context.Context, template *domaincompute.LaunchTemplate) (*domaincompute.LaunchTemplate, error) {
	if err := template.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsTemplate := awsmapper.FromDomainLaunchTemplate(template)
	if err := awsTemplate.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsTemplateOutput, err := a.awsService.CreateLaunchTemplate(ctx, awsTemplate)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTemplate := awsmapper.ToDomainLaunchTemplateFromOutput(awsTemplateOutput)
	// Preserve region from input
	domainTemplate.Region = template.Region
	return domainTemplate, nil
}

func (a *AWSComputeAdapter) GetLaunchTemplate(ctx context.Context, id string) (*domaincompute.LaunchTemplate, error) {
	awsTemplateOutput, err := a.awsService.GetLaunchTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainLaunchTemplateFromOutput(awsTemplateOutput), nil
}

func (a *AWSComputeAdapter) UpdateLaunchTemplate(ctx context.Context, id string, template *domaincompute.LaunchTemplate) (*domaincompute.LaunchTemplate, error) {
	if err := template.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsTemplate := awsmapper.FromDomainLaunchTemplate(template)
	if err := awsTemplate.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsTemplateOutput, err := a.awsService.UpdateLaunchTemplate(ctx, id, awsTemplate)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTemplate := awsmapper.ToDomainLaunchTemplateFromOutput(awsTemplateOutput)
	// Preserve region from input
	domainTemplate.Region = template.Region
	return domainTemplate, nil
}

func (a *AWSComputeAdapter) DeleteLaunchTemplate(ctx context.Context, id string) error {
	if err := a.awsService.DeleteLaunchTemplate(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListLaunchTemplates(ctx context.Context, filters map[string]string) ([]*domaincompute.LaunchTemplate, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsTemplateOutputs, err := a.awsService.ListLaunchTemplates(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTemplates := make([]*domaincompute.LaunchTemplate, len(awsTemplateOutputs))
	for i, awsTemplateOutput := range awsTemplateOutputs {
		domainTemplates[i] = awsmapper.ToDomainLaunchTemplateFromOutput(awsTemplateOutput)
	}

	return domainTemplates, nil
}

func (a *AWSComputeAdapter) GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*domaincompute.LaunchTemplate, error) {
	awsVersion, err := a.awsService.GetLaunchTemplateVersion(ctx, id, version)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTemplate := awsmapper.ToDomainLaunchTemplate(awsVersion.TemplateData)
	if domainTemplate != nil {
		domainTemplate.ID = awsVersion.TemplateID
		domainTemplate.Version = &awsVersion.VersionNumber
	}

	return domainTemplate, nil
}

func (a *AWSComputeAdapter) ListLaunchTemplateVersions(ctx context.Context, id string) ([]*domaincompute.LaunchTemplateVersion, error) {
	awsVersions, err := a.awsService.ListLaunchTemplateVersions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVersions := make([]*domaincompute.LaunchTemplateVersion, len(awsVersions))
	for i := range awsVersions {
		domainVersions[i] = awsmapper.ToDomainLaunchTemplateVersion(awsVersions[i])
	}

	return domainVersions, nil
}
