package compute

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// InstanceOutput represents the output data for a compute instance after creation/update
// This includes cloud-generated identifiers and runtime state
type InstanceOutput struct {
	// Core identifiers
	ID               string
	ARN              *string
	Name             string
	Region           string
	AvailabilityZone *string

	// Compute Configuration
	InstanceType string
	AMI          string

	// Networking
	SubnetID         string
	SecurityGroupIDs []string
	PrivateIP        *string
	PublicIP         *string
	PrivateDNS       *string
	PublicDNS        *string
	VPCID            *string

	// Access & Permissions
	KeyName            *string
	IAMInstanceProfile *string

	// State
	State     InstanceState
	CreatedAt *time.Time
}

// LoadBalancerOutput represents the output data for a load balancer after creation/update
type LoadBalancerOutput struct {
	// Core identifiers
	ID     string
	ARN    *string
	Name   string
	Region string

	// Configuration
	Type             LoadBalancerType
	Internal         bool
	SecurityGroupIDs []string
	SubnetIDs        []string

	// Output fields (cloud-generated)
	DNSName *string // Auto-generated DNS name
	ZoneID  *string // Route53 hosted zone ID

	// State
	State     LoadBalancerState
	CreatedAt *time.Time
}

// TargetGroupOutput represents the output data for a target group after creation/update
type TargetGroupOutput struct {
	// Core identifiers
	ID    string
	ARN   *string
	Name  string
	VPCID string

	// Configuration
	Port        int
	Protocol    TargetGroupProtocol
	TargetType  TargetType
	HealthCheck HealthCheckConfig

	// State
	State     TargetGroupState
	CreatedAt *time.Time
}

// ListenerOutput represents the output data for a listener after creation/update
type ListenerOutput struct {
	// Core identifiers
	ID              string
	ARN             *string
	LoadBalancerARN string

	// Configuration
	Port          int
	Protocol      ListenerProtocol
	DefaultAction ListenerAction
	Rules         []ListenerRule

	// CreatedAt timestamp
	CreatedAt *time.Time
}

// LaunchTemplateOutput represents the output data for a launch template after creation/update
type LaunchTemplateOutput struct {
	// Core identifiers
	ID         string
	ARN        *string
	Name       string
	Region     string
	NamePrefix *string

	// Version Management
	DefaultVersion *int
	LatestVersion  *int

	// CreatedAt timestamp
	CreatedAt *time.Time
	CreatedBy *string
}

// AutoScalingGroupOutput represents the output data for an auto scaling group after creation/update
type AutoScalingGroupOutput struct {
	// Core identifiers
	ID         string
	ARN        *string
	Name       string
	Region     string
	NamePrefix *string

	// Capacity Configuration
	MinSize         int
	MaxSize         int
	DesiredCapacity *int

	// Location Configuration
	VPCZoneIdentifier []string

	// Launch Configuration
	LaunchTemplate *LaunchTemplateSpecification

	// Health Check Configuration
	HealthCheckType        AutoScalingGroupHealthCheckType
	HealthCheckGracePeriod *int

	// Load Balancer Integration
	TargetGroupARNs []string

	// State
	State       AutoScalingGroupState
	CreatedTime *string // ISO 8601 timestamp
}

// LambdaFunctionOutput represents the output data for a Lambda function after creation/update
type LambdaFunctionOutput struct {
	// Core identifiers
	FunctionName string
	ARN          *string
	InvokeARN    *string
	QualifiedARN *string
	Region       string

	// Code Source
	S3Bucket        *string
	S3Key           *string
	S3ObjectVersion *string
	PackageType     *string
	ImageURI        *string

	// Runtime Configuration
	Runtime *string
	Handler *string

	// Configuration
	MemorySize  *int
	Timeout     *int
	Environment map[string]string
	Layers      []string
	VPCConfig   *LambdaVPCConfig

	// Output fields (cloud-generated)
	Version      *string
	LastModified *string
	CodeSize     *int64
	CodeSHA256   *string

	// CreatedAt timestamp
	CreatedAt *time.Time
}

// TargetGroupAttachmentOutput represents the output data for a target group attachment
type TargetGroupAttachmentOutput struct {
	TargetGroupARN   string
	TargetID         string
	Port             *int
	AvailabilityZone *string
	State            *string
}

// ToInstanceOutput converts an Instance domain model to InstanceOutput
func ToInstanceOutput(instance *Instance) *InstanceOutput {
	if instance == nil {
		return nil
	}
	return &InstanceOutput{
		ID:                 instance.ID,
		ARN:                instance.ARN,
		Name:               instance.Name,
		Region:             instance.Region,
		AvailabilityZone:   instance.AvailabilityZone,
		InstanceType:       instance.InstanceType,
		AMI:                instance.AMI,
		SubnetID:           instance.SubnetID,
		SecurityGroupIDs:   instance.SecurityGroupIDs,
		PrivateIP:          instance.PrivateIP,
		PublicIP:           instance.PublicIP,
		KeyName:            instance.KeyName,
		IAMInstanceProfile: instance.IAMInstanceProfile,
		State:              instance.State,
	}
}

// ToLoadBalancerOutput converts a LoadBalancer domain model to LoadBalancerOutput
func ToLoadBalancerOutput(lb *LoadBalancer) *LoadBalancerOutput {
	if lb == nil {
		return nil
	}
	return &LoadBalancerOutput{
		ID:               lb.ID,
		ARN:              lb.ARN,
		Name:             lb.Name,
		Region:           lb.Region,
		Type:             lb.Type,
		Internal:         lb.Internal,
		SecurityGroupIDs: lb.SecurityGroupIDs,
		SubnetIDs:        lb.SubnetIDs,
		DNSName:          lb.DNSName,
		ZoneID:           lb.ZoneID,
		State:            lb.State,
	}
}

// ToTargetGroupOutput converts a TargetGroup domain model to TargetGroupOutput
func ToTargetGroupOutput(tg *TargetGroup) *TargetGroupOutput {
	if tg == nil {
		return nil
	}
	return &TargetGroupOutput{
		ID:          tg.ID,
		ARN:         tg.ARN,
		Name:        tg.Name,
		VPCID:       tg.VPCID,
		Port:        tg.Port,
		Protocol:    tg.Protocol,
		TargetType:  tg.TargetType,
		HealthCheck: tg.HealthCheck,
		State:       tg.State,
	}
}

// ToListenerOutput converts a Listener domain model to ListenerOutput
func ToListenerOutput(listener *Listener) *ListenerOutput {
	if listener == nil {
		return nil
	}
	return &ListenerOutput{
		ID:              listener.ID,
		ARN:             listener.ARN,
		LoadBalancerARN: listener.LoadBalancerARN,
		Port:            listener.Port,
		Protocol:        listener.Protocol,
		DefaultAction:   listener.DefaultAction,
		Rules:           listener.Rules,
	}
}

// ToLaunchTemplateOutput converts a LaunchTemplate domain model to LaunchTemplateOutput
func ToLaunchTemplateOutput(template *LaunchTemplate) *LaunchTemplateOutput {
	if template == nil {
		return nil
	}
	return &LaunchTemplateOutput{
		ID:             template.ID,
		ARN:            template.ARN,
		Name:           template.Name,
		Region:         template.Region,
		NamePrefix:     template.NamePrefix,
		DefaultVersion: template.Version,
		LatestVersion:  template.LatestVersion,
	}
}

// ToAutoScalingGroupOutput converts an AutoScalingGroup domain model to AutoScalingGroupOutput
func ToAutoScalingGroupOutput(asg *AutoScalingGroup) *AutoScalingGroupOutput {
	if asg == nil {
		return nil
	}
	return &AutoScalingGroupOutput{
		ID:                     asg.ID,
		ARN:                    asg.ARN,
		Name:                   asg.Name,
		Region:                 asg.Region,
		NamePrefix:             asg.NamePrefix,
		MinSize:                asg.MinSize,
		MaxSize:                asg.MaxSize,
		DesiredCapacity:        asg.DesiredCapacity,
		VPCZoneIdentifier:      asg.VPCZoneIdentifier,
		LaunchTemplate:         asg.LaunchTemplate,
		HealthCheckType:        asg.HealthCheckType,
		HealthCheckGracePeriod: asg.HealthCheckGracePeriod,
		TargetGroupARNs:        asg.TargetGroupARNs,
		State:                  asg.State,
		CreatedTime:            asg.CreatedTime,
	}
}

// ToLambdaFunctionOutput converts a LambdaFunction domain model to LambdaFunctionOutput
func ToLambdaFunctionOutput(function *LambdaFunction) *LambdaFunctionOutput {
	if function == nil {
		return nil
	}
	return &LambdaFunctionOutput{
		FunctionName:    function.FunctionName,
		ARN:             function.ARN,
		InvokeARN:       function.InvokeARN,
		QualifiedARN:    function.QualifiedARN,
		Region:          function.Region,
		S3Bucket:        function.S3Bucket,
		S3Key:           function.S3Key,
		S3ObjectVersion: function.S3ObjectVersion,
		PackageType:     function.PackageType,
		ImageURI:        function.ImageURI,
		Runtime:         function.Runtime,
		Handler:         function.Handler,
		MemorySize:      function.MemorySize,
		Timeout:         function.Timeout,
		Environment:     function.Environment,
		Layers:          function.Layers,
		VPCConfig:       function.VPCConfig,
		Version:         function.Version,
		LastModified:    function.LastModified,
		CodeSize:        function.CodeSize,
		CodeSHA256:      function.CodeSHA256,
	}
}

// ToResourceOutput converts an InstanceOutput to the generic ResourceOutput
func (io *InstanceOutput) ToResourceOutput() *resource.ResourceOutput {
	return &resource.ResourceOutput{
		ID:       io.ID,
		ARN:      io.ARN,
		Name:     io.Name,
		Region:   io.Region,
		State:    (*string)(&io.State),
		Provider: resource.AWS, // This would need to be passed in or determined
	}
}
