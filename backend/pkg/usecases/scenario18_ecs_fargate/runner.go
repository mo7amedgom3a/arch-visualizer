package scenario18_ecs_fargate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/compute"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/containers"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/iam"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/terraform"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	tfgen "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/generator"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func Run(ctx context.Context) error {
	computeMode := os.Getenv("COMPUTE_MODE")
	if computeMode == "" {
		computeMode = "FARGATE"
	}

	fmt.Println("==================================================================================")
	fmt.Printf("SCENARIO 18: ECS %s with ALB\n", computeMode)
	fmt.Println("==================================================================================")

	// Construct an ECS architecture
	arch := architecture.NewArchitecture()
	arch.Region = "us-east-1"
	arch.Provider = resource.AWS

	// Resource IDs
	vpcID := "vpc-ecs"
	subnetPublic1ID := "subnet-public-1"
	subnetPublic2ID := "subnet-public-2"
	subnetPrivate1ID := "subnet-private-1"
	subnetPrivate2ID := "subnet-private-2"
	sgALBID := "sg-alb"
	sgECSID := "sg-ecs"
	albID := "alb-ecs"
	tgID := "tg-ecs"
	listenerID := "listener-ecs"
	ecsClusterID := "ecs-cluster"
	ecsTaskDefID := "ecs-task-def"
	ecsServiceID := "ecs-service"
	igwID := "igw-ecs"

	// EC2 Mode specific IDs
	ecrRepoID := "ecr-repo"
	launchTemplateID := "ecs-lt"
	asgID := "ecs-asg"
	capacityProviderID := "ecs-cp"
	clusterCapacityProvidersID := "ecs-cluster-cps"
	iamRoleExecutionID := "ecs-execution-role"
	iamInstanceProfileID := "ecs-instance-profile"

	// 1. VPC
	vpc := &resource.Resource{
		ID:       vpcID,
		Name:     "ecs-vpc",
		Type:     resource.ResourceType{Name: "VPC", ID: "VPC", Category: "Networking", Kind: "VPC"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"cidr":               "10.0.0.0/16",
			"enableDnsHostnames": true,
			"enableDnsSupport":   true,
		},
	}
	arch.Resources = append(arch.Resources, vpc)

	// 2. Internet Gateway
	igw := &resource.Resource{
		ID:        igwID,
		Name:      "ecs-igw",
		Type:      resource.ResourceType{Name: "InternetGateway", ID: "InternetGateway", Category: "Networking", Kind: "Gateway"},
		Provider:  resource.AWS,
		ParentID:  &vpcID,
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, igw)

	// 3. Public Subnets
	subnetPublic1 := &resource.Resource{
		ID:       subnetPublic1ID,
		Name:     "public-subnet-1",
		Type:     resource.ResourceType{Name: "Subnet", ID: "Subnet", Category: "Networking", Kind: "Subnet"},
		ParentID: &vpcID,
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"cidr":                    "10.0.1.0/24",
			"availabilityZoneId":      "us-east-1a",
			"map_public_ip_on_launch": true,
		},
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, subnetPublic1)

	subnetPublic2 := &resource.Resource{
		ID:       subnetPublic2ID,
		Name:     "public-subnet-2",
		Type:     resource.ResourceType{Name: "Subnet", ID: "Subnet", Category: "Networking", Kind: "Subnet"},
		ParentID: &vpcID,
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"cidr":                    "10.0.2.0/24",
			"availabilityZoneId":      "us-east-1b",
			"map_public_ip_on_launch": true,
		},
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, subnetPublic2)

	// 4. Private Subnets
	subnetPrivate1 := &resource.Resource{
		ID:       subnetPrivate1ID,
		Name:     "private-subnet-1",
		Type:     resource.ResourceType{Name: "Subnet", ID: "Subnet", Category: "Networking", Kind: "Subnet"},
		ParentID: &vpcID,
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"cidr":               "10.0.3.0/24",
			"availabilityZoneId": "us-east-1a",
		},
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, subnetPrivate1)

	subnetPrivate2 := &resource.Resource{
		ID:       subnetPrivate2ID,
		Name:     "private-subnet-2",
		Type:     resource.ResourceType{Name: "Subnet", ID: "Subnet", Category: "Networking", Kind: "Subnet"},
		ParentID: &vpcID,
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"cidr":               "10.0.4.0/24",
			"availabilityZoneId": "us-east-1b",
		},
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, subnetPrivate2)

	// 5. ALB Security Group
	sgALB := &resource.Resource{
		ID:       sgALBID,
		Name:     "alb-sg",
		Type:     resource.ResourceType{Name: "SecurityGroup", ID: "SecurityGroup", Category: "Networking", Kind: "SecurityGroup"},
		ParentID: &vpcID,
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"description": "Allow HTTP from internet",
			"ingressRules": []interface{}{
				map[string]interface{}{
					"protocol":    "tcp",
					"fromPort":    80,
					"toPort":      80,
					"cidr":        "0.0.0.0/0",
					"description": "HTTP",
				},
			},
			"egressRules": []interface{}{
				map[string]interface{}{
					"protocol": "-1",
					"fromPort": 0,
					"toPort":   0,
					"cidr":     "0.0.0.0/0",
				},
			},
		},
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, sgALB)

	// 6. ECS Security Group
	sgECS := &resource.Resource{
		ID:       sgECSID,
		Name:     "ecs-tasks-sg",
		Type:     resource.ResourceType{Name: "SecurityGroup", ID: "SecurityGroup", Category: "Networking", Kind: "SecurityGroup"},
		ParentID: &vpcID,
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"description": "Allow HTTP from ALB",
			"ingressRules": []interface{}{
				map[string]interface{}{
					"protocol":              "tcp",
					"fromPort":              80,
					"toPort":                80,
					"sourceSecurityGroupId": sgALBID,
					"description":           "HTTP from ALB",
				},
			},
			"egressRules": []interface{}{
				map[string]interface{}{
					"protocol": "-1",
					"fromPort": 0,
					"toPort":   0,
					"cidr":     "0.0.0.0/0",
				},
			},
		},
		DependsOn: []string{vpcID, sgALBID},
	}
	arch.Resources = append(arch.Resources, sgECS)

	// 7. Target Group
	targetType := "ip"
	if computeMode == "EC2" {
		targetType = "instance"
	}

	tg := &resource.Resource{
		ID:       tgID,
		Name:     "ecs-tg",
		Type:     resource.ResourceType{Name: "TargetGroup", ID: "TargetGroup", Category: "Compute", Kind: "TargetGroup"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"vpcId":      vpcID,
			"port":       80,
			"protocol":   "HTTP",
			"targetType": targetType,
			"healthCheck": map[string]interface{}{
				"path":     "/health",
				"protocol": "HTTP",
			},
		},
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, tg)

	// 8. Application Load Balancer
	alb := &resource.Resource{
		ID:       albID,
		Name:     "ecs-alb",
		Type:     resource.ResourceType{Name: "LoadBalancer", ID: "LoadBalancer", Category: "Compute", Kind: "LoadBalancer"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"load_balancer_type": "application",
			"internal":           false,
			"subnets": []interface{}{
				map[string]interface{}{"subnetId": subnetPublic1ID},
				map[string]interface{}{"subnetId": subnetPublic2ID},
			},
			"securityGroupIds": []string{sgALBID},
		},
		DependsOn: []string{subnetPublic1ID, subnetPublic2ID, sgALBID},
	}
	arch.Resources = append(arch.Resources, alb)

	// 9. Listener
	listener := &resource.Resource{
		ID:       listenerID,
		Name:     "ecs-listener",
		Type:     resource.ResourceType{Name: "Listener", ID: "Listener", Category: "Compute", Kind: "Listener"},
		ParentID: &albID,
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"port":              80,
			"protocol":          "HTTP",
			"defaultActionType": "forward",
			"targetGroupId":     tgID,
		},
		DependsOn: []string{albID, tgID},
	}
	arch.Resources = append(arch.Resources, listener)

	// 10. ECS Cluster
	ecsCluster := &resource.Resource{
		ID:       ecsClusterID,
		Name:     "main-cluster",
		Type:     resource.ResourceType{Name: "ECSCluster", ID: "ECSCluster", Category: "Containers", Kind: "Container"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"containerInsightsEnabled": true,
			"executeCommandEnabled":    true,
		},
	}
	arch.Resources = append(arch.Resources, ecsCluster)

	if computeMode == "EC2" {
		// IAM Role for EC2 Instances
		roleEC2ID := "role-ec2"
		roleEC2 := &resource.Resource{
			ID:       roleEC2ID,
			Name:     "ecs-instance-role",
			Type:     resource.ResourceType{Name: "IAMRole", ID: "IAMRole", Category: "IAM", Kind: "Role"},
			Provider: resource.AWS,
			Metadata: map[string]interface{}{
				"assumeRolePolicyDocument": `{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Action": "sts:AssumeRole",
							"Principal": {
								"Service": "ec2.amazonaws.com"
							},
							"Effect": "Allow",
							"Sid": ""
						}
					]
				}`,
				"managedPolicyArns": []string{
					"arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role",
				},
			},
		}
		arch.Resources = append(arch.Resources, roleEC2)

		// IAM Instance Profile
		instanceProfile := &resource.Resource{
			ID:       iamInstanceProfileID,
			Name:     "ecs-instance-profile",
			Type:     resource.ResourceType{Name: "IAMInstanceProfile", ID: "IAMInstanceProfile", Category: "IAM", Kind: "InstanceProfile"},
			Provider: resource.AWS,
			Metadata: map[string]interface{}{
				"role": "ecs-instance-role",
			},
			DependsOn: []string{roleEC2ID},
		}
		arch.Resources = append(arch.Resources, instanceProfile)

		// IAM Role for ECS Task Execution
		roleExecution := &resource.Resource{
			ID:       iamRoleExecutionID,
			Name:     "ecs-task-execution-role",
			Type:     resource.ResourceType{Name: "IAMRole", ID: "IAMRole", Category: "IAM", Kind: "Role"},
			Provider: resource.AWS,
			Metadata: map[string]interface{}{
				"assumeRolePolicyDocument": `{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Action": "sts:AssumeRole",
							"Principal": {
								"Service": "ecs-tasks.amazonaws.com"
							},
							"Effect": "Allow",
							"Sid": ""
						}
					]
				}`,
				"managedPolicyArns": []string{
					"arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
				},
			},
		}
		arch.Resources = append(arch.Resources, roleExecution)

		// ECR Repository
		ecrRepo := &resource.Resource{
			ID:       ecrRepoID,
			Name:     "app-repo",
			Type:     resource.ResourceType{Name: "ECRRepository", ID: "ECRRepository", Category: "Containers", Kind: "Container"},
			Provider: resource.AWS,
			Metadata: map[string]interface{}{
				"image_tag_mutability": "MUTABLE",
				"scan_on_push":         true,
				"force_delete":         true,
			},
		}
		arch.Resources = append(arch.Resources, ecrRepo)

		// Launch Template
		lt := &resource.Resource{
			ID:       launchTemplateID,
			Name:     "ecs-template",
			Type:     resource.ResourceType{Name: "LaunchTemplate", ID: "LaunchTemplate", Category: "Compute", Kind: "LaunchTemplate"},
			Provider: resource.AWS,
			Metadata: map[string]interface{}{
				"imageId":            "ami-ecs-optimized",
				"instanceType":       "t3.micro",
				"userData":           fmt.Sprintf("#!/bin/bash\necho ECS_CLUSTER=%s >> /etc/ecs/ecs.config", "main-cluster"),
				"iamInstanceProfile": iamInstanceProfileID,
				"securityGroupIds":   []string{sgECSID},
			},
		}
		arch.Resources = append(arch.Resources, lt)

		// Auto Scaling Group
		asg := &resource.Resource{
			ID:       asgID,
			Name:     "ecs-asg",
			Type:     resource.ResourceType{Name: "AutoScalingGroup", ID: "AutoScalingGroup", Category: "Compute", Kind: "AutoScalingGroup"},
			Provider: resource.AWS,
			Metadata: map[string]interface{}{
				"minSize":           1,
				"maxSize":           4,
				"desiredCapacity":   2,
				"vpcZoneIdentifier": []string{subnetPrivate1ID, subnetPrivate2ID},
				"launchTemplateId":  launchTemplateID,
				"tags": []interface{}{
					map[string]interface{}{"key": "Name", "value": "ecs-node", "propagateAtLaunch": true},
				},
			},
			DependsOn: []string{launchTemplateID, subnetPrivate1ID, subnetPrivate2ID},
		}
		arch.Resources = append(arch.Resources, asg)

		// Capacity Provider
		cp := &resource.Resource{
			ID:       capacityProviderID,
			Name:     "ecs-ec2-provider",
			Type:     resource.ResourceType{Name: "ECSCapacityProvider", ID: "ECSCapacityProvider", Category: "Containers", Kind: "Container"},
			Provider: resource.AWS,
			Metadata: map[string]interface{}{
				"auto_scaling_group_arn":         asgID, // Mapper should resolve this ID to ARN
				"managed_termination_protection": "DISABLED",
				"managed_scaling_enabled":        true,
				"target_capacity":                80,
			},
			DependsOn: []string{asgID},
		}
		arch.Resources = append(arch.Resources, cp)

		// Cluster Capacity Providers Attachment
		clusterCP := &resource.Resource{
			ID:       clusterCapacityProvidersID,
			Name:     "main-cluster", // Same name as cluster to associate
			Type:     resource.ResourceType{Name: "ECSClusterCapacityProviders", ID: "ECSClusterCapacityProviders", Category: "Containers", Kind: "Container"},
			Provider: resource.AWS,
			Metadata: map[string]interface{}{
				"cluster_name":       "main-cluster",
				"capacity_providers": []string{capacityProviderID}, // Helper will need to resolve name
				"default_capacity_provider_strategy": []interface{}{
					map[string]interface{}{
						"capacity_provider": capacityProviderID,
						"weight":            1,
						"base":              1,
					},
				},
			},
			DependsOn: []string{ecsClusterID, capacityProviderID},
		}
		arch.Resources = append(arch.Resources, clusterCP)
	}

	// 11. ECS Task Definition
	networkMode := "awsvpc"
	requiresCompatibilities := []string{"FARGATE"}

	if computeMode == "EC2" {
		networkMode = "bridge"
		requiresCompatibilities = []string{"EC2"}
	}

	ecsTaskDef := &resource.Resource{
		ID:       ecsTaskDefID,
		Name:     "web-app",
		Type:     resource.ResourceType{Name: "ECSTaskDefinition", ID: "ECSTaskDefinition", Category: "Containers", Kind: "Container"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"family":                  "web-app",
			"networkMode":             networkMode,
			"requiresCompatibilities": requiresCompatibilities,
			"cpu":                     "256",
			"memory":                  "512",
			"executionRoleArn":        iamRoleExecutionID,
			"containerDefinitions": []interface{}{
				map[string]interface{}{
					"name":      "web-container",
					"image":     "nginx:latest", // Simplified, would use ECR URL in real scenario
					"cpu":       256,
					"memory":    512,
					"essential": true,
					"portMappings": []interface{}{
						map[string]interface{}{
							"containerPort": 80,
							"hostPort":      80, // Dynamic port mapping for bridge mode? Or static 80 for logic simplicity
							"protocol":      "tcp",
						},
					},
				},
			},
		},
	}
	arch.Resources = append(arch.Resources, ecsTaskDef)

	// 12. ECS Service
	serviceLaunchType := "FARGATE"
	var capacityProviderStrategy []interface{}

	if computeMode == "EC2" {
		serviceLaunchType = "" // Launch type must be empty when using capacity provider strategy
		capacityProviderStrategy = []interface{}{
			map[string]interface{}{
				"capacity_provider": "ecs-ec2-provider", // Name of the CP resource
				"weight":            1,
			},
		}
	}

	ecsServiceMetadata := map[string]interface{}{
		"clusterName":    "main-cluster",
		"taskDefinition": ecsTaskDefID,
		"desiredCount":   2,
		"loadBalancer": map[string]interface{}{
			"targetGroupArn": tgID,
			"containerName":  "web-container",
			"containerPort":  80,
		},
		"deploymentCircuitBreaker": map[string]interface{}{
			"enable":   true,
			"rollback": true,
		},
	}

	if computeMode == "FARGATE" {
		ecsServiceMetadata["launchType"] = serviceLaunchType
		ecsServiceMetadata["networkConfiguration"] = map[string]interface{}{
			"subnets":        []string{subnetPrivate1ID, subnetPrivate2ID},
			"securityGroups": []string{sgECSID},
			"assignPublicIp": false,
		}
	} else {
		// EC2 Mode
		ecsServiceMetadata["capacity_provider_strategy"] = capacityProviderStrategy
		// Network config is optional for bridge mode, but usually required for awsvpc
		// With bridge mode, no network config needed or different
	}

	ecsService := &resource.Resource{
		ID:        ecsServiceID,
		Name:      "web-service",
		Type:      resource.ResourceType{Name: "ECSService", ID: "ECSService", Category: "Containers", Kind: "Container"},
		Provider:  resource.AWS,
		Metadata:  ecsServiceMetadata,
		DependsOn: []string{ecsClusterID, ecsTaskDefID, tgID, listenerID},
	}

	if computeMode == "EC2" {
		ecsService.DependsOn = append(ecsService.DependsOn, clusterCapacityProvidersID)
	}

	arch.Resources = append(arch.Resources, ecsService)

	// Sort resources
	graph := architecture.NewGraph(arch)
	sortedResources, err := graph.GetSortedResources()
	if err != nil {
		return fmt.Errorf("getting sorted resources: %w", err)
	}

	fmt.Printf("✓ Created architecture with %d resources\n", len(arch.Resources))

	// Generate Terraform
	mapperRegistry := tfmapper.NewRegistry()
	awsMapper := terraform.New()
	if err := mapperRegistry.Register(awsMapper); err != nil {
		return fmt.Errorf("register aws terraform mapper: %w", err)
	}

	engine := tfgen.NewEngine(mapperRegistry)
	output, err := engine.Generate(ctx, arch, sortedResources)
	if err != nil {
		return fmt.Errorf("terraform engine generate: %w", err)
	}

	// Print output
	fmt.Println("\n==================================================================================")
	fmt.Println("Generated Terraform Files:")
	fmt.Println("==================================================================================")
	for _, f := range output.Files {
		fmt.Printf("\n--- %s ---\n", f.Path)
		fmt.Println(f.Content)
	}

	// Write to file for inspection
	outDir := filepath.Join("terraform_output", "scenario18_ecs_fargate")
	if err := writeTerraformOutput(outDir, output); err != nil {
		return err
	}
	fmt.Printf("\n✓ Files written to %s\n", outDir)

	return nil
}

func writeTerraformOutput(dir string, out *iac.Output) error {
	if out == nil {
		return fmt.Errorf("nil terraform output")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output dir %s: %w", dir, err)
	}

	for _, f := range out.Files {
		target := filepath.Join(dir, f.Path)
		if err := os.WriteFile(target, []byte(f.Content), 0o644); err != nil {
			return fmt.Errorf("write file %s: %w", target, err)
		}
	}
	return nil
}
