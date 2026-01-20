package common

import (
	"fmt"
	"time"

	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
)

// MockIDGenerator generates consistent mock IDs for resources
type MockIDGenerator struct {
	counter int
}

// NewMockIDGenerator creates a new mock ID generator
func NewMockIDGenerator() *MockIDGenerator {
	return &MockIDGenerator{counter: 1}
}

// GenerateVPCID generates a mock VPC ID
func (g *MockIDGenerator) GenerateVPCID() string {
	id := fmt.Sprintf("vpc-%08x", g.counter)
	g.counter++
	return id
}

// GenerateSubnetID generates a mock subnet ID
func (g *MockIDGenerator) GenerateSubnetID() string {
	id := fmt.Sprintf("subnet-%08x", g.counter)
	g.counter++
	return id
}

// GenerateInternetGatewayID generates a mock Internet Gateway ID
func (g *MockIDGenerator) GenerateInternetGatewayID() string {
	id := fmt.Sprintf("igw-%08x", g.counter)
	g.counter++
	return id
}

// GenerateRouteTableID generates a mock route table ID
func (g *MockIDGenerator) GenerateRouteTableID() string {
	id := fmt.Sprintf("rtb-%08x", g.counter)
	g.counter++
	return id
}

// GenerateSecurityGroupID generates a mock security group ID
func (g *MockIDGenerator) GenerateSecurityGroupID() string {
	id := fmt.Sprintf("sg-%08x", g.counter)
	g.counter++
	return id
}

// GenerateNATGatewayID generates a mock NAT Gateway ID
func (g *MockIDGenerator) GenerateNATGatewayID() string {
	id := fmt.Sprintf("nat-%08x", g.counter)
	g.counter++
	return id
}

// GenerateElasticIPID generates a mock Elastic IP allocation ID
func (g *MockIDGenerator) GenerateElasticIPID() string {
	id := fmt.Sprintf("eipalloc-%08x", g.counter)
	g.counter++
	return id
}

// GenerateInstanceID generates a mock EC2 instance ID
func (g *MockIDGenerator) GenerateInstanceID() string {
	id := fmt.Sprintf("i-%08x", g.counter)
	g.counter++
	return id
}

// GenerateLoadBalancerARN generates a mock load balancer ARN
func (g *MockIDGenerator) GenerateLoadBalancerARN(region, accountID, name string) string {
	return fmt.Sprintf("arn:aws:elasticloadbalancing:%s:%s:loadbalancer/app/%s/%s", region, accountID, name, generateRandomHex(16))
}

// GenerateTargetGroupARN generates a mock target group ARN
func (g *MockIDGenerator) GenerateTargetGroupARN(region, accountID, name string) string {
	return fmt.Sprintf("arn:aws:elasticloadbalancing:%s:%s:targetgroup/%s/%s", region, accountID, name, generateRandomHex(16))
}

// GenerateLaunchTemplateID generates a mock launch template ID
func (g *MockIDGenerator) GenerateLaunchTemplateID() string {
	id := fmt.Sprintf("lt-%08x", g.counter)
	g.counter++
	return id
}

// GenerateASGARN generates a mock Auto Scaling Group ARN
func (g *MockIDGenerator) GenerateASGARN(region, accountID, name string) string {
	return fmt.Sprintf("arn:aws:autoscaling:%s:%s:autoScalingGroup:uuid:autoScalingGroupName/%s", region, accountID, name)
}

// GenerateIAMRoleARN generates a mock IAM role ARN
func (g *MockIDGenerator) GenerateIAMRoleARN(region, accountID, roleName string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
}

// GenerateIAMInstanceProfileARN generates a mock IAM instance profile ARN
func (g *MockIDGenerator) GenerateIAMInstanceProfileARN(accountID, profileName string) string {
	return fmt.Sprintf("arn:aws:iam::%s:instance-profile/%s", accountID, profileName)
}

// GenerateARN generates a generic ARN
func (g *MockIDGenerator) GenerateARN(service, region, accountID, resourceType, resourceID string) string {
	return fmt.Sprintf("arn:aws:%s:%s:%s:%s/%s", service, region, accountID, resourceType, resourceID)
}

// generateRandomHex generates a random hex string
func generateRandomHex(length int) string {
	chars := "0123456789abcdef"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
	}
	return string(result)
}

// CreateMockVPC creates a mock VPC domain model
func CreateMockVPC(region, name, cidr string, gen *MockIDGenerator) *domainnetworking.VPC {
	vpcID := gen.GenerateVPCID()
	arn := gen.GenerateARN("ec2", region, "123456789012", "vpc", vpcID)
	return &domainnetworking.VPC{
		ID:                 vpcID,
		ARN:                &arn,
		Name:               name,
		Region:             region,
		CIDR:               cidr,
		EnableDNS:          true,
		EnableDNSHostnames: true,
	}
}

// CreateMockSubnet creates a mock subnet domain model
func CreateMockSubnet(vpcID, name, cidr, az string, gen *MockIDGenerator) *domainnetworking.Subnet {
	subnetID := gen.GenerateSubnetID()
	azPtr := &az
	return &domainnetworking.Subnet{
		ID:               subnetID,
		Name:             name,
		VPCID:            vpcID,
		CIDR:             cidr,
		AvailabilityZone: azPtr,
	}
}

// CreateMockInternetGateway creates a mock Internet Gateway domain model
func CreateMockInternetGateway(vpcID, name string, gen *MockIDGenerator) *domainnetworking.InternetGateway {
	igwID := gen.GenerateInternetGatewayID()
	arn := gen.GenerateARN("ec2", "us-east-1", "123456789012", "internet-gateway", igwID)
	return &domainnetworking.InternetGateway{
		ID:    igwID,
		ARN:   &arn,
		Name:  name,
		VPCID: vpcID,
	}
}

// CreateMockRouteTable creates a mock route table domain model
func CreateMockRouteTable(vpcID, name string, gen *MockIDGenerator) *domainnetworking.RouteTable {
	rtID := gen.GenerateRouteTableID()
	return &domainnetworking.RouteTable{
		ID:    rtID,
		Name:  name,
		VPCID: vpcID,
	}
}

// CreateMockSecurityGroup creates a mock security group domain model
func CreateMockSecurityGroup(vpcID, name, description string, gen *MockIDGenerator) *domainnetworking.SecurityGroup {
	sgID := gen.GenerateSecurityGroupID()
	return &domainnetworking.SecurityGroup{
		ID:          sgID,
		Name:        name,
		Description: description,
		VPCID:       vpcID,
	}
}

// CreateMockNATGateway creates a mock NAT Gateway domain model
func CreateMockNATGateway(subnetID, name string, gen *MockIDGenerator) *domainnetworking.NATGateway {
	natID := gen.GenerateNATGatewayID()
	eipID := gen.GenerateElasticIPID()
	arn := gen.GenerateARN("ec2", "us-east-1", "123456789012", "nat-gateway", natID)
	return &domainnetworking.NATGateway{
		ID:           natID,
		ARN:          &arn,
		Name:         name,
		SubnetID:     subnetID,
		AllocationID: &eipID,
	}
}

// CreateMockEC2Instance creates a mock EC2 instance domain model
func CreateMockEC2Instance(name, instanceType, subnetID, sgID, region, az string, gen *MockIDGenerator) *domaincompute.Instance {
	instanceID := gen.GenerateInstanceID()
	arn := gen.GenerateARN("ec2", region, "123456789012", "instance", instanceID)
	azPtr := &az
	privateIP := "10.0.1.100"
	return &domaincompute.Instance{
		ID:               instanceID,
		ARN:              &arn,
		Name:             name,
		Region:           region,
		AvailabilityZone: azPtr,
		InstanceType:     instanceType,
		AMI:              "ami-0c55b159cbfafe1f0",
		SubnetID:         subnetID,
		SecurityGroupIDs: []string{sgID},
		PrivateIP:        &privateIP,
		State:            domaincompute.InstanceStateRunning,
	}
}

// CreateMockLoadBalancer creates a mock load balancer domain model
func CreateMockLoadBalancer(name, lbType string, subnetIDs, sgIDs []string, region string, gen *MockIDGenerator) *domaincompute.LoadBalancer {
	arn := gen.GenerateLoadBalancerARN(region, "123456789012", name)
	dnsName := fmt.Sprintf("%s-%s.elb.amazonaws.com", name, region)
	zoneID := "Z35SXDOTRQ7X7K"
	return &domaincompute.LoadBalancer{
		ID:               name,
		ARN:              &arn,
		Name:             name,
		Region:           region,
		Type:             domaincompute.LoadBalancerType(lbType),
		Internal:         false,
		SecurityGroupIDs: sgIDs,
		SubnetIDs:        subnetIDs,
		DNSName:          &dnsName,
		ZoneID:           &zoneID,
		State:            domaincompute.LoadBalancerStateActive,
	}
}

// CreateMockTargetGroup creates a mock target group domain model
func CreateMockTargetGroup(name, vpcID, protocol string, port int, region string, gen *MockIDGenerator) *domaincompute.TargetGroup {
	arn := gen.GenerateTargetGroupARN(region, "123456789012", name)
	healthPath := "/health"
	healthPort := fmt.Sprintf("%d", port)
	healthProtocol := protocol
	return &domaincompute.TargetGroup{
		ID:         name,
		ARN:        &arn,
		Name:       name,
		VPCID:      vpcID,
		Protocol:   domaincompute.TargetGroupProtocol(protocol),
		Port:       port,
		TargetType: domaincompute.TargetTypeInstance,
		HealthCheck: domaincompute.HealthCheckConfig{
			Path:     &healthPath,
			Protocol: &healthProtocol,
			Port:     &healthPort,
		},
		State: domaincompute.TargetGroupStateActive,
	}
}

// CreateMockLaunchTemplate creates a mock launch template domain model
func CreateMockLaunchTemplate(name, instanceType, imageID string, sgIDs []string, region string, gen *MockIDGenerator) *domaincompute.LaunchTemplate {
	ltID := gen.GenerateLaunchTemplateID()
	arn := gen.GenerateARN("ec2", region, "123456789012", "launch-template", ltID)
	return &domaincompute.LaunchTemplate{
		ID:               ltID,
		ARN:              &arn,
		Name:             name,
		Region:           region,
		ImageID:          imageID,
		InstanceType:     instanceType,
		SecurityGroupIDs: sgIDs,
	}
}

// CreateMockAutoScalingGroup creates a mock Auto Scaling Group domain model
func CreateMockAutoScalingGroup(name string, minSize, maxSize, desiredCapacity int, subnetIDs []string, launchTemplateID, region string, gen *MockIDGenerator) *domaincompute.AutoScalingGroup {
	arn := gen.GenerateASGARN(region, "123456789012", name)
	version := "$Latest"
	gracePeriod := 300
	createdTime := time.Now().Format("2006-01-02T15:04:05Z07:00")
	return &domaincompute.AutoScalingGroup{
		ID:                name,
		ARN:               &arn,
		Name:              name,
		Region:            region,
		MinSize:           minSize,
		MaxSize:           maxSize,
		DesiredCapacity:   &desiredCapacity,
		VPCZoneIdentifier: subnetIDs,
		LaunchTemplate: &domaincompute.LaunchTemplateSpecification{
			ID:      launchTemplateID,
			Version: &version,
		},
		HealthCheckType:        domaincompute.AutoScalingGroupHealthCheckTypeEC2,
		HealthCheckGracePeriod: &gracePeriod,
		TargetGroupARNs:        []string{},
		Tags:                   []domaincompute.Tag{},
		State:                  domaincompute.AutoScalingGroupStateActive,
		CreatedTime:            &createdTime,
	}
}

// CreateMockIAMRole creates a mock IAM role domain model
func CreateMockIAMRole(name, description string, gen *MockIDGenerator) *domainiam.Role {
	roleID := gen.GenerateIAMRoleARN("us-east-1", "123456789012", name)
	path := "/"
	assumeRolePolicy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "ec2.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`
	return &domainiam.Role{
		ID:               name,
		ARN:              &roleID,
		Name:             name,
		Description:      &description,
		Path:             &path,
		AssumeRolePolicy: assumeRolePolicy,
	}
}

// CreateMockIAMInstanceProfile creates a mock IAM instance profile domain model
func CreateMockIAMInstanceProfile(name, roleName string, gen *MockIDGenerator) *domainiam.InstanceProfile {
	profileARN := gen.GenerateIAMInstanceProfileARN("123456789012", name)
	path := "/"
	return &domainiam.InstanceProfile{
		ID:       name,
		ARN:      &profileARN,
		Name:     name,
		Path:     &path,
		RoleName: &roleName,
	}
}

// GetDefaultAccountID returns the default mock account ID
func GetDefaultAccountID() string {
	return "123456789012"
}

// GetDefaultAvailabilityZones returns default availability zones for a region
func GetDefaultAvailabilityZones(region string) []string {
	azMap := map[string][]string{
		"us-east-1":      {"us-east-1a", "us-east-1b", "us-east-1c"},
		"us-west-2":      {"us-west-2a", "us-west-2b", "us-west-2c"},
		"eu-west-1":      {"eu-west-1a", "eu-west-1b", "eu-west-1c"},
		"ap-southeast-1": {"ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"},
	}
	if azs, ok := azMap[region]; ok {
		return azs
	}
	// Default fallback
	return []string{fmt.Sprintf("%sa", region), fmt.Sprintf("%sb", region), fmt.Sprintf("%sc", region)}
}
