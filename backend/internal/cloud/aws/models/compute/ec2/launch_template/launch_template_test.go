package launch_template

import (
	"strings"
	"testing"
)

func validLaunchTemplate() *LaunchTemplate {
	namePrefix := "test-template"
	return &LaunchTemplate{
		NamePrefix:          &namePrefix,
		ImageID:             "ami-0c55b159cbfafe1f0",
		InstanceType:        "t3.micro",
		VpcSecurityGroupIds: []string{"sg-123"},
	}
}

func TestLaunchTemplateValidate_Success(t *testing.T) {
	lt := validLaunchTemplate()

	if err := lt.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLaunchTemplateValidate_NameOrPrefixRequired(t *testing.T) {
	lt := validLaunchTemplate()
	lt.NamePrefix = nil
	lt.Name = nil

	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for missing name and name_prefix")
	}
}

func TestLaunchTemplateValidate_ImageIDValidation(t *testing.T) {
	lt := validLaunchTemplate()
	lt.ImageID = "bad-image"
	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for invalid image_id prefix")
	}

	lt.ImageID = "ami-123"
	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for invalid image_id length")
	}
}

func TestLaunchTemplateValidate_InstanceTypeValidation(t *testing.T) {
	lt := validLaunchTemplate()
	lt.InstanceType = "not-valid"

	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for invalid instance type")
	}
}

func TestLaunchTemplateValidate_SecurityGroupValidation(t *testing.T) {
	lt := validLaunchTemplate()
	lt.VpcSecurityGroupIds = nil
	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for missing security groups")
	}

	lt = validLaunchTemplate()
	lt.VpcSecurityGroupIds = []string{"bad-123"}
	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for invalid security group id")
	}
}

func TestLaunchTemplateValidate_UserDataSize(t *testing.T) {
	lt := validLaunchTemplate()
	tooLarge := strings.Repeat("a", 16385)
	lt.UserData = &tooLarge

	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for oversized user data")
	}
}

func TestLaunchTemplateValidate_RootVolumeValidation(t *testing.T) {
	lt := validLaunchTemplate()
	root := "bad-vol"
	lt.RootVolumeID = &root

	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for invalid root volume id")
	}
}

func TestLaunchTemplateValidate_AdditionalVolumeValidation(t *testing.T) {
	lt := validLaunchTemplate()
	lt.AdditionalVolumeIDs = []string{"", "vol-123"}
	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for empty additional volume id")
	}

	lt = validLaunchTemplate()
	lt.AdditionalVolumeIDs = []string{"bad-123"}
	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for invalid additional volume id")
	}
}

func TestLaunchTemplateValidate_IAMInstanceProfileValidation(t *testing.T) {
	lt := validLaunchTemplate()
	lt.IAMInstanceProfile = &IAMInstanceProfile{}

	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for invalid IAM instance profile")
	}
}

func TestLaunchTemplateValidate_MetadataOptionsValidation(t *testing.T) {
	lt := validLaunchTemplate()
	endpoint := "invalid"
	lt.MetadataOptions = &MetadataOptions{HTTPEndpoint: &endpoint}

	if err := lt.Validate(); err == nil {
		t.Fatal("expected error for invalid metadata options")
	}
}
