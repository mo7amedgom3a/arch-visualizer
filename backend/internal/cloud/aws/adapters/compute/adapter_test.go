package compute

import (
	"context"
	"errors"
	"testing"
	"time"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsec2 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
	awsec2outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/outputs"
	awslttemplate "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template"
	awslttemplateoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template/outputs"
	awsInstanceTypes "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// mockAWSComputeService is a mock implementation of AWSComputeService for testing
type mockAWSComputeService struct {
	instance        *awsec2.Instance
	launchTemplate  *awslttemplate.LaunchTemplate
	createError     error
	getError        error
	instanceTypeInfo *awsInstanceTypes.InstanceTypeInfo
}

// Ensure mockAWSComputeService implements AWSComputeService
var _ awsservice.AWSComputeService = (*mockAWSComputeService)(nil)

// Helper function to convert Instance input to output
func instanceToOutput(instance *awsec2.Instance) *awsec2outputs.InstanceOutput {
	if instance == nil {
		return nil
	}
	return &awsec2outputs.InstanceOutput{
		ID:                "i-mock-1234567890abcdef0",
		ARN:               "arn:aws:ec2:us-east-1:123456789012:instance/i-mock-1234567890abcdef0",
		Name:              instance.Name,
		Region:            "us-east-1",
		AvailabilityZone:  "us-east-1a",
		InstanceType:      instance.InstanceType,
		AMI:               instance.AMI,
		SubnetID:          instance.SubnetID,
		SecurityGroupIDs:  instance.VpcSecurityGroupIds,
		PrivateIP:         "10.0.1.100",
		PublicIP:          nil,
		KeyName:           instance.KeyName,
		IAMInstanceProfile: instance.IAMInstanceProfile,
		State:             "running",
		CreationTime:      time.Now(),
		Tags:              instance.Tags,
	}
}

// Helper function to convert LaunchTemplate input to output
func launchTemplateToOutput(template *awslttemplate.LaunchTemplate) *awslttemplateoutputs.LaunchTemplateOutput {
	if template == nil {
		return nil
	}
	name := "test-template"
	if template.Name != nil {
		name = *template.Name
	} else if template.NamePrefix != nil {
		name = *template.NamePrefix + "-12345"
	}
	return &awslttemplateoutputs.LaunchTemplateOutput{
		ID:             "lt-mock-1234567890abcdef0",
		ARN:            "arn:aws:ec2:us-east-1:123456789012:launch-template/lt-mock-1234567890abcdef0",
		Name:           name,
		DefaultVersion: 1,
		LatestVersion:  1,
		CreateTime:     time.Now(),
		CreatedBy:      "arn:aws:iam::123456789012:user/test",
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: name},
		},
	}
}

// Instance Operations

func (m *mockAWSComputeService) CreateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.instance = instance
	return instanceToOutput(instance), nil
}

func (m *mockAWSComputeService) GetInstance(ctx context.Context, id string) (*awsec2outputs.InstanceOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return instanceToOutput(m.instance), nil
}

func (m *mockAWSComputeService) UpdateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	m.instance = instance
	return instanceToOutput(instance), nil
}

func (m *mockAWSComputeService) DeleteInstance(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSComputeService) ListInstances(ctx context.Context, filters map[string][]string) ([]*awsec2outputs.InstanceOutput, error) {
	if m.instance != nil {
		return []*awsec2outputs.InstanceOutput{instanceToOutput(m.instance)}, nil
	}
	return []*awsec2outputs.InstanceOutput{}, nil
}

// Instance Lifecycle Operations

func (m *mockAWSComputeService) StartInstance(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSComputeService) StopInstance(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSComputeService) RebootInstance(ctx context.Context, id string) error {
	return nil
}

// Instance Type Operations

func (m *mockAWSComputeService) GetInstanceTypeInfo(ctx context.Context, instanceType string, region string) (*awsInstanceTypes.InstanceTypeInfo, error) {
	if m.instanceTypeInfo != nil {
		return m.instanceTypeInfo, nil
	}
	return &awsInstanceTypes.InstanceTypeInfo{
		Name: instanceType,
		Category:     awsInstanceTypes.CategoryGeneralPurpose,
		VCPU:         2,
		MemoryGiB:       4.0,
		MaxNetworkGbps:     5,
		StorageType:     "EBS only",
	}, nil
}

	func (m *mockAWSComputeService) ListInstanceTypesByCategory(ctx context.Context, category awsInstanceTypes.InstanceCategory, region string) ([]*awsInstanceTypes.InstanceTypeInfo, error) {
	return []*awsInstanceTypes.InstanceTypeInfo{
		{
			Name: "t3.micro",
			Category:     category,
			VCPU:         2,
			MemoryGiB:    1.0,
			MaxNetworkGbps: 5,
			StorageType:   "EBS",
		},
	}, nil
}

// Launch Template Operations

func (m *mockAWSComputeService) CreateLaunchTemplate(ctx context.Context, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.launchTemplate = template
	return launchTemplateToOutput(template), nil
}

func (m *mockAWSComputeService) GetLaunchTemplate(ctx context.Context, id string) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return launchTemplateToOutput(m.launchTemplate), nil
}

func (m *mockAWSComputeService) UpdateLaunchTemplate(ctx context.Context, id string, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	m.launchTemplate = template
	return launchTemplateToOutput(template), nil
}

func (m *mockAWSComputeService) DeleteLaunchTemplate(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSComputeService) ListLaunchTemplates(ctx context.Context, filters map[string][]string) ([]*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if m.launchTemplate != nil {
		return []*awslttemplateoutputs.LaunchTemplateOutput{launchTemplateToOutput(m.launchTemplate)}, nil
	}
	return []*awslttemplateoutputs.LaunchTemplateOutput{}, nil
}

func (m *mockAWSComputeService) GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*awslttemplate.LaunchTemplateVersion, error) {
	return &awslttemplate.LaunchTemplateVersion{
		TemplateID:    id,
		VersionNumber: version,
		IsDefault:     version == 1,
		CreateTime:    time.Now(),
		CreatedBy:     stringPtr("arn:aws:iam::123456789012:user/test"),
		TemplateData:  m.launchTemplate,
	}, nil
}

func (m *mockAWSComputeService) ListLaunchTemplateVersions(ctx context.Context, id string) ([]*awslttemplate.LaunchTemplateVersion, error) {
	return []*awslttemplate.LaunchTemplateVersion{
		{
			TemplateID:    id,
			VersionNumber: 1,
			IsDefault:     true,
			CreateTime:    time.Now(),
			CreatedBy:     stringPtr("arn:aws:iam::123456789012:user/test"),
			TemplateData:  m.launchTemplate,
		},
	}, nil
}


// Tests

func TestAWSComputeAdapter_CreateInstance(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	domainInstance := &domaincompute.Instance{
		Name:         "test-instance",
		Region:       "us-east-1",
		InstanceType: "t3.micro",
		AMI:          "ami-0c55b159cbfafe1f0",
		SubnetID:     "subnet-123",
		SecurityGroupIDs: []string{"sg-123"},
	}

	ctx := context.Background()
	createdInstance, err := adapter.CreateInstance(ctx, domainInstance)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdInstance == nil {
		t.Fatal("Expected created instance, got nil")
	}

	if createdInstance.Name != domainInstance.Name {
		t.Errorf("Expected name %s, got %s", domainInstance.Name, createdInstance.Name)
	}

	if createdInstance.ID == "" {
		t.Error("Expected instance ID to be populated")
	}
}

func TestAWSComputeAdapter_CreateInstance_ValidationError(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	invalidInstance := &domaincompute.Instance{
		Name:   "", // Invalid: empty name
		Region: "us-east-1",
		AMI:    "ami-0c55b159cbfafe1f0",
	}

	ctx := context.Background()
	_, err := adapter.CreateInstance(ctx, invalidInstance)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestAWSComputeAdapter_GetInstance(t *testing.T) {
	mockService := &mockAWSComputeService{
		instance: &awsec2.Instance{
			Name:         "test-instance",
			InstanceType: "t3.micro",
			AMI:          "ami-0c55b159cbfafe1f0",
			SubnetID:     "subnet-123",
			VpcSecurityGroupIds: []string{"sg-123"},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	instance, err := adapter.GetInstance(ctx, "i-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance, got nil")
	}

	if instance.Name != "test-instance" {
		t.Errorf("Expected name test-instance, got %s", instance.Name)
	}
}

func TestAWSComputeAdapter_ListInstances(t *testing.T) {
	mockService := &mockAWSComputeService{
		instance: &awsec2.Instance{
			Name:         "test-instance",
			InstanceType: "t3.micro",
			AMI:          "ami-0c55b159cbfafe1f0",
			SubnetID:     "subnet-123",
			VpcSecurityGroupIds: []string{"sg-123"},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	instances, err := adapter.ListInstances(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(instances) != 1 {
		t.Errorf("Expected 1 instance, got %d", len(instances))
	}

	if instances[0].Name != "test-instance" {
		t.Errorf("Expected name test-instance, got %s", instances[0].Name)
	}
}

func TestAWSComputeAdapter_CreateLaunchTemplate(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	namePrefix := "test-template"
	domainTemplate := &domaincompute.LaunchTemplate{
		Name:         "test-template",
		Region:       "us-east-1",
		NamePrefix:   &namePrefix,
		ImageID:      "ami-0c55b159cbfafe1f0",
		InstanceType: "t3.micro",
		SecurityGroupIDs: []string{"sg-123"},
	}

	ctx := context.Background()
	createdTemplate, err := adapter.CreateLaunchTemplate(ctx, domainTemplate)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdTemplate == nil {
		t.Fatal("Expected created template, got nil")
	}

	if createdTemplate.ID == "" {
		t.Error("Expected template ID to be populated")
	}
}

func TestAWSComputeAdapter_ErrorHandling(t *testing.T) {
	mockService := &mockAWSComputeService{
		getError: errors.New("aws service error"),
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	_, err := adapter.GetInstance(ctx, "i-123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify error is wrapped
	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}
