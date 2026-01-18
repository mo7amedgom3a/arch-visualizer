package compute

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsec2 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
	awslttemplate "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template"
	awslttemplateoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template/outputs"
	awsec2outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/outputs"
	awsInstanceTypes "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// realisticAWSComputeService is a realistic implementation that returns proper output models
type realisticAWSComputeService struct{}

var _ awsservice.AWSComputeService = (*realisticAWSComputeService)(nil)

// Instance Operations

func (s *realisticAWSComputeService) CreateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	// Simulate realistic AWS instance creation
	return &awsec2outputs.InstanceOutput{
		ID:                 "i-0a1b2c3d4e5f6g7h8",
		ARN:                "arn:aws:ec2:us-east-1:123456789012:instance/i-0a1b2c3d4e5f6g7h8",
		Name:               instance.Name,
		Region:             "us-east-1",
		AvailabilityZone:   "us-east-1a",
		InstanceType:       instance.InstanceType,
		AMI:                instance.AMI,
		SubnetID:           instance.SubnetID,
		SecurityGroupIDs:   instance.VpcSecurityGroupIds,
		PrivateIP:          "10.0.1.100",
		PublicIP:           stringPtr("54.123.45.67"),
		KeyName:            instance.KeyName,
		IAMInstanceProfile: instance.IAMInstanceProfile,
		State:              "running",
		CreationTime:       time.Now(),
		Tags:               instance.Tags,
	}, nil
}

func (s *realisticAWSComputeService) GetInstance(ctx context.Context, id string) (*awsec2outputs.InstanceOutput, error) {
	return &awsec2outputs.InstanceOutput{
		ID:                 id,
		ARN:                fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:instance/%s", id),
		Name:               "test-instance",
		Region:             "us-east-1",
		AvailabilityZone:   "us-east-1a",
		InstanceType:       "t3.micro",
		AMI:                "ami-0c55b159cbfafe1f0",
		SubnetID:           "subnet-123",
		SecurityGroupIDs:   []string{"sg-123"},
		PrivateIP:          "10.0.1.100",
		PublicIP:           stringPtr("54.123.45.67"),
		KeyName:            stringPtr("my-key"),
		IAMInstanceProfile: nil,
		State:              "running",
		CreationTime:       time.Now(),
		Tags:               []configs.Tag{},
	}, nil
}

func (s *realisticAWSComputeService) UpdateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	return s.CreateInstance(context.Background(), instance)
}

func (s *realisticAWSComputeService) DeleteInstance(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSComputeService) ListInstances(ctx context.Context, filters map[string][]string) ([]*awsec2outputs.InstanceOutput, error) {
	return []*awsec2outputs.InstanceOutput{
		{
			ID:                 "i-0a1b2c3d4e5f6g7h8",
			ARN:                "arn:aws:ec2:us-east-1:123456789012:instance/i-0a1b2c3d4e5f6g7h8",
			Name:               "test-instance",
			Region:             "us-east-1",
			AvailabilityZone:   "us-east-1a",
			InstanceType:       "t3.micro",
			AMI:                "ami-0c55b159cbfafe1f0",
			SubnetID:           "subnet-123",
			SecurityGroupIDs:   []string{"sg-123"},
			PrivateIP:          "10.0.1.100",
			PublicIP:           stringPtr("54.123.45.67"),
			KeyName:            stringPtr("my-key"),
			IAMInstanceProfile: nil,
			State:              "running",
			CreationTime:       time.Now(),
			Tags:               []configs.Tag{},
		},
	}, nil
}

// Instance Lifecycle Operations

func (s *realisticAWSComputeService) StartInstance(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSComputeService) StopInstance(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSComputeService) RebootInstance(ctx context.Context, id string) error {
	return nil
}

// Instance Type Operations

func (s *realisticAWSComputeService) GetInstanceTypeInfo(ctx context.Context, instanceType string, region string) (*awsInstanceTypes.InstanceTypeInfo, error) {
	return &awsInstanceTypes.InstanceTypeInfo{
		Name:           instanceType,
		Category:       awsInstanceTypes.CategoryGeneralPurpose,
		VCPU:           2,
		MemoryGiB:      4.0,
		MaxNetworkGbps: 5,
		StorageType:    "EBS only",
	}, nil
}

func (s *realisticAWSComputeService) ListInstanceTypesByCategory(ctx context.Context, category awsInstanceTypes.InstanceCategory, region string) ([]*awsInstanceTypes.InstanceTypeInfo, error) {
	return []*awsInstanceTypes.InstanceTypeInfo{
		{
			Name:           "t3.micro",
			Category:       category,
			VCPU:           2,
			MemoryGiB:      1.0,
			MaxNetworkGbps: 5,
			StorageType:    "EBS",
		},
	}, nil
}

// Launch Template Operations

func (s *realisticAWSComputeService) CreateLaunchTemplate(ctx context.Context, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	name := "test-template"
	if template.Name != nil {
		name = *template.Name
	} else if template.NamePrefix != nil {
		name = *template.NamePrefix + "-12345"
	}
	return &awslttemplateoutputs.LaunchTemplateOutput{
		ID:             "lt-0a1b2c3d4e5f6g7h8",
		ARN:            "arn:aws:ec2:us-east-1:123456789012:launch-template/lt-0a1b2c3d4e5f6g7h8",
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
	}, nil
}

func (s *realisticAWSComputeService) GetLaunchTemplate(ctx context.Context, id string) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return &awslttemplateoutputs.LaunchTemplateOutput{
		ID:             id,
		ARN:            fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:launch-template/%s", id),
		Name:           "test-template",
		DefaultVersion: 1,
		LatestVersion:  1,
		CreateTime:     time.Now(),
		CreatedBy:      "arn:aws:iam::123456789012:user/test",
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: "test-template"},
		},
	}, nil
}

func (s *realisticAWSComputeService) UpdateLaunchTemplate(ctx context.Context, id string, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return s.CreateLaunchTemplate(context.Background(), template)
}

func (s *realisticAWSComputeService) DeleteLaunchTemplate(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSComputeService) ListLaunchTemplates(ctx context.Context, filters map[string][]string) ([]*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return []*awslttemplateoutputs.LaunchTemplateOutput{
		{
			ID:             "lt-0a1b2c3d4e5f6g7h8",
			ARN:            "arn:aws:ec2:us-east-1:123456789012:launch-template/lt-0a1b2c3d4e5f6g7h8",
			Name:           "test-template",
			DefaultVersion: 1,
			LatestVersion:  1,
			CreateTime:     time.Now(),
			CreatedBy:      "arn:aws:iam::123456789012:user/test",
			Tags: []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				{Key: "Name", Value: "test-template"},
			},
		},
	}, nil
}

func (s *realisticAWSComputeService) GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*awslttemplate.LaunchTemplateVersion, error) {
	templateData := &awslttemplate.LaunchTemplate{
		Name:                stringPtr("test-template"),
		ImageID:             "ami-0c55b159cbfafe1f0",
		InstanceType:        "t3.micro",
		VpcSecurityGroupIds: []string{"sg-123"},
	}
	return &awslttemplate.LaunchTemplateVersion{
		TemplateID:    id,
		VersionNumber: version,
		IsDefault:     version == 1,
		CreateTime:    time.Now(),
		CreatedBy:     stringPtr("arn:aws:iam::123456789012:user/test"),
		TemplateData:  templateData,
	}, nil
}

func (s *realisticAWSComputeService) ListLaunchTemplateVersions(ctx context.Context, id string) ([]*awslttemplate.LaunchTemplateVersion, error) {
	return []*awslttemplate.LaunchTemplateVersion{
		{
			TemplateID:    id,
			VersionNumber: 1,
			IsDefault:     true,
			CreateTime:    time.Now(),
			CreatedBy:     stringPtr("arn:aws:iam::123456789012:user/test"),
			TemplateData:  nil,
		},
	}, nil
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

// Integration Tests

func TestAWSComputeAdapter_OutputIntegration_CreateInstance(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	domainInstance := &domaincompute.Instance{
		Name:             "integration-test-instance",
		Region:           "us-east-1",
		InstanceType:     "t3.micro",
		AMI:              "ami-0c55b159cbfafe1f0",
		SubnetID:         "subnet-123",
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

	// Verify AWS-generated identifiers are populated
	if createdInstance.ID == "" {
		t.Error("Expected instance ID to be populated")
	}

	if createdInstance.ID != "i-0a1b2c3d4e5f6g7h8" {
		t.Errorf("Expected instance ID i-0a1b2c3d4e5f6g7h8, got %s", createdInstance.ID)
	}

	if createdInstance.ARN == nil {
		t.Error("Expected instance ARN to be populated")
	}

	if createdInstance.ARN != nil && *createdInstance.ARN == "" {
		t.Error("Expected instance ARN to be non-empty")
	}

	// Verify state is populated
	if createdInstance.State != domaincompute.InstanceStateRunning {
		t.Errorf("Expected state running, got %s", createdInstance.State)
	}

	// Verify IP addresses are populated
	if createdInstance.PrivateIP == nil {
		t.Error("Expected private IP to be populated")
	}

	if createdInstance.PrivateIP != nil && *createdInstance.PrivateIP != "10.0.1.100" {
		t.Errorf("Expected private IP 10.0.1.100, got %s", *createdInstance.PrivateIP)
	}
}

func TestAWSComputeAdapter_OutputIntegration_GetInstance(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	ctx := context.Background()
	instance, err := adapter.GetInstance(ctx, "i-0a1b2c3d4e5f6g7h8")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance, got nil")
	}

	// Verify all output fields are populated
	if instance.ID != "i-0a1b2c3d4e5f6g7h8" {
		t.Errorf("Expected ID i-0a1b2c3d4e5f6g7h8, got %s", instance.ID)
	}

	if instance.ARN == nil {
		t.Error("Expected ARN to be populated")
	}

	if instance.State == "" {
		t.Error("Expected state to be populated")
	}

	if instance.PrivateIP == nil {
		t.Error("Expected private IP to be populated")
	}
}

func TestAWSComputeAdapter_OutputIntegration_ListInstances(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	ctx := context.Background()
	instances, err := adapter.ListInstances(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(instances) == 0 {
		t.Fatal("Expected at least one instance, got none")
	}

	// Verify first instance has all output fields
	instance := instances[0]
	if instance.ID == "" {
		t.Error("Expected instance ID to be populated")
	}

	if instance.ARN == nil {
		t.Error("Expected instance ARN to be populated")
	}

	if instance.State == "" {
		t.Error("Expected instance state to be populated")
	}
}

func TestAWSComputeAdapter_OutputIntegration_CreateLaunchTemplate(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	namePrefix := "integration-template"
	domainTemplate := &domaincompute.LaunchTemplate{
		Name:             "integration-template",
		Region:           "us-east-1",
		NamePrefix:       &namePrefix,
		ImageID:          "ami-0c55b159cbfafe1f0",
		InstanceType:     "t3.micro",
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

	// Verify AWS-generated identifiers are populated
	if createdTemplate.ID == "" {
		t.Error("Expected template ID to be populated")
	}

	if createdTemplate.ARN == nil {
		t.Error("Expected template ARN to be populated")
	}

	// Verify version information is populated
	if createdTemplate.Version == nil {
		t.Error("Expected version to be populated")
	}

	if createdTemplate.LatestVersion == nil {
		t.Error("Expected latest version to be populated")
	}
}

func TestAWSComputeAdapter_OutputIntegration_GetLaunchTemplate(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	ctx := context.Background()
	template, err := adapter.GetLaunchTemplate(ctx, "lt-0a1b2c3d4e5f6g7h8")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if template == nil {
		t.Fatal("Expected template, got nil")
	}

	// Verify all output fields are populated
	if template.ID != "lt-0a1b2c3d4e5f6g7h8" {
		t.Errorf("Expected ID lt-0a1b2c3d4e5f6g7h8, got %s", template.ID)
	}

	if template.ARN == nil {
		t.Error("Expected ARN to be populated")
	}

	if template.Version == nil {
		t.Error("Expected version to be populated")
	}
}

func TestAWSComputeAdapter_OutputIntegration_ListLaunchTemplates(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	ctx := context.Background()
	templates, err := adapter.ListLaunchTemplates(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(templates) == 0 {
		t.Fatal("Expected at least one template, got none")
	}

	// Verify first template has all output fields
	template := templates[0]
	if template.ID == "" {
		t.Error("Expected template ID to be populated")
	}

	if template.ARN == nil {
		t.Error("Expected template ARN to be populated")
	}

	if template.Version == nil {
		t.Error("Expected template version to be populated")
	}
}

func TestAWSComputeAdapter_OutputIntegration_GetLaunchTemplateVersion(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	ctx := context.Background()
	template, err := adapter.GetLaunchTemplateVersion(ctx, "lt-0a1b2c3d4e5f6g7h8", 1)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if template == nil {
		t.Fatal("Expected template, got nil")
	}

	// Verify version information
	if template.Version == nil {
		t.Error("Expected version to be populated")
	}

	if template.Version != nil && *template.Version != 1 {
		t.Errorf("Expected version 1, got %d", *template.Version)
	}
}

func TestAWSComputeAdapter_OutputIntegration_ListLaunchTemplateVersions(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	ctx := context.Background()
	versions, err := adapter.ListLaunchTemplateVersions(ctx, "lt-0a1b2c3d4e5f6g7h8")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(versions) == 0 {
		t.Fatal("Expected at least one version, got none")
	}

	// Verify first version has all fields
	version := versions[0]
	if version.TemplateID == "" {
		t.Error("Expected template ID to be populated")
	}

	if version.VersionNumber == 0 {
		t.Error("Expected version number to be populated")
	}
}

func TestAWSComputeAdapter_OutputIntegration_InstanceLifecycle(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService)

	ctx := context.Background()

	// Test StartInstance
	err := adapter.StartInstance(ctx, "i-0a1b2c3d4e5f6g7h8")
	if err != nil {
		t.Fatalf("Expected no error starting instance, got: %v", err)
	}

	// Test StopInstance
	err = adapter.StopInstance(ctx, "i-0a1b2c3d4e5f6g7h8")
	if err != nil {
		t.Fatalf("Expected no error stopping instance, got: %v", err)
	}

	// Test RebootInstance
	err = adapter.RebootInstance(ctx, "i-0a1b2c3d4e5f6g7h8")
	if err != nil {
		t.Fatalf("Expected no error rebooting instance, got: %v", err)
	}
}

func TestAWSComputeAdapter_OutputIntegration_GetInstanceTypeInfo(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService).(*AWSComputeAdapter)

	ctx := context.Background()
	info, err := adapter.GetInstanceTypeInfo(ctx, "t3.micro", "us-east-1")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if info == nil {
		t.Fatal("Expected instance type info, got nil")
	}

	if info.Name != "t3.micro" {
		t.Errorf("Expected instance type t3.micro, got %s", info.Name)
	}

	if info.VCPU == 0 {
		t.Error("Expected VCPU to be populated")
	}

	if info.MemoryGiB == 0 {
		t.Error("Expected memory to be populated")
	}
}

func TestAWSComputeAdapter_OutputIntegration_ListInstanceTypesByCategory(t *testing.T) {
	realisticService := &realisticAWSComputeService{}
	adapter := NewAWSComputeAdapter(realisticService).(*AWSComputeAdapter)

	ctx := context.Background()
	types, err := adapter.ListInstanceTypesByCategory(ctx, awsInstanceTypes.CategoryGeneralPurpose, "us-east-1")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(types) == 0 {
		t.Fatal("Expected at least one instance type, got none")
	}

	// Verify first type has all fields
	instanceType := types[0]
	if instanceType.Name == "" {
		t.Error("Expected instance type to be populated")
	}

	if instanceType.Category != awsInstanceTypes.CategoryGeneralPurpose {
		t.Error("Expected category to be populated")
	}
}
