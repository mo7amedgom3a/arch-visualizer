package launch_template

import "testing"

func TestMetadataOptionsValidate_Success(t *testing.T) {
	endpoint := "enabled"
	tokens := "required"
	hopLimit := 2
	mo := &MetadataOptions{
		HTTPEndpoint:            &endpoint,
		HTTPTokens:              &tokens,
		HTTPPutResponseHopLimit: &hopLimit,
	}

	if err := mo.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMetadataOptionsValidate_InvalidValues(t *testing.T) {
	invalidEndpoint := "on"
	mo := &MetadataOptions{HTTPEndpoint: &invalidEndpoint}
	if err := mo.Validate(); err == nil {
		t.Fatal("expected error for invalid http_endpoint")
	}

	invalidTokens := "maybe"
	mo = &MetadataOptions{HTTPTokens: &invalidTokens}
	if err := mo.Validate(); err == nil {
		t.Fatal("expected error for invalid http_tokens")
	}

	tooHigh := 100
	mo = &MetadataOptions{HTTPPutResponseHopLimit: &tooHigh}
	if err := mo.Validate(); err == nil {
		t.Fatal("expected error for invalid http_put_response_hop_limit")
	}
}

func TestIAMInstanceProfileValidate_Success(t *testing.T) {
	name := "profile-name"
	iip := &IAMInstanceProfile{Name: &name}
	if err := iip.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	arn := "arn:aws:iam::123456789012:instance-profile/profile-name"
	iip = &IAMInstanceProfile{ARN: &arn}
	if err := iip.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestIAMInstanceProfileValidate_Invalid(t *testing.T) {
	iip := &IAMInstanceProfile{}
	if err := iip.Validate(); err == nil {
		t.Fatal("expected error for missing name and arn")
	}

	invalidARN := "arn:aws:ec2::123"
	iip = &IAMInstanceProfile{ARN: &invalidARN}
	if err := iip.Validate(); err == nil {
		t.Fatal("expected error for invalid arn")
	}
}
