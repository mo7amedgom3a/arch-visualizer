package scenario17_alb_asg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/terraform"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	tfgen "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/generator"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func Run(ctx context.Context) error {
	// Construct an architecture manually
	arch := architecture.NewArchitecture()
	arch.Region = "us-east-1"
	arch.Provider = resource.AWS

	// Helpers values
	vpcID := "vpc-1"
	subnetID := "subnet-1"
	sgID := "sg-1"
	ltID := "lt-1"
	tgID := "tg-1"
	asgID := "asg-1"
	albID := "alb-1"
	listenerID := "listener-1"

	// 1. VPC
	vpc := &resource.Resource{
		ID:       vpcID,
		Name:     "MainVPC",
		Type:     resource.ResourceType{Name: "VPC", ID: "VPC", Category: "Networking", Kind: "VPC"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"cidr": "10.0.0.0/16",
		},
	}
	arch.Resources = append(arch.Resources, vpc)

	// 2. Subnet
	subnet := &resource.Resource{
		ID:       subnetID,
		Name:     "PublicSubnet",
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
	arch.Resources = append(arch.Resources, subnet)

	// 3. Security Group
	sg := &resource.Resource{
		ID:       sgID,
		Name:     "WebSG",
		Type:     resource.ResourceType{Name: "SecurityGroup", ID: "SecurityGroup", Category: "Networking", Kind: "SecurityGroup"},
		ParentID: &vpcID,
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"description": "Allow HTTP",
			"ingressRules": []interface{}{
				map[string]interface{}{
					"protocol":    "tcp",
					"fromPort":    80,
					"toPort":      80,
					"cidr":        "0.0.0.0/0",
					"description": "HTTP",
				},
			},
		},
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, sg)

	// 4. Launch Template
	lt := &resource.Resource{
		ID:       ltID,
		Name:     "WebLaunchTemplate",
		Type:     resource.ResourceType{Name: "LaunchTemplate", ID: "LaunchTemplate", Category: "Compute", Kind: "LaunchTemplate"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"imageId":          "ami-0123456789",
			"instanceType":     "t3.micro",
			"securityGroupIds": []string{sgID},
		},
		DependsOn: []string{sgID},
	}
	arch.Resources = append(arch.Resources, lt)

	// 5. Target Group
	tg := &resource.Resource{
		ID:       tgID,
		Name:     "WebTargetGroup",
		Type:     resource.ResourceType{Name: "TargetGroup", ID: "TargetGroup", Category: "Compute", Kind: "TargetGroup"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"vpcId":      vpcID,
			"port":       80,
			"protocol":   "HTTP",
			"targetType": "instance",
			"healthCheck": map[string]interface{}{
				"path": "/",
			},
		},
		DependsOn: []string{vpcID},
	}
	arch.Resources = append(arch.Resources, tg)

	// 6. Auto Scaling Group
	asg := &resource.Resource{
		ID:       asgID,
		Name:     "WebASG",
		Type:     resource.ResourceType{Name: "AutoScalingGroup", ID: "AutoScalingGroup", Category: "Compute", Kind: "AutoScalingGroup"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"minSize":          1,
			"maxSize":          3,
			"desiredCapacity":  2,
			"launchTemplateId": ltID,
			"subnets": []interface{}{
				map[string]interface{}{"subnetId": subnetID},
			},
			"targetGroupIds": []string{tgID},
		},
		DependsOn: []string{ltID, subnetID, tgID},
	}
	arch.Resources = append(arch.Resources, asg)

	// 7. Load Balancer (ALB)
	alb := &resource.Resource{
		ID:       albID,
		Name:     "WebALB",
		Type:     resource.ResourceType{Name: "LoadBalancer", ID: "LoadBalancer", Category: "Compute", Kind: "LoadBalancer"},
		Provider: resource.AWS,
		Metadata: map[string]interface{}{
			"load_balancer_type": "application",
			"subnets": []interface{}{
				map[string]interface{}{"subnetId": subnetID},
			},
			"securityGroupIds": []string{sgID},
		},
		DependsOn: []string{subnetID, sgID},
	}
	arch.Resources = append(arch.Resources, alb)

	// 8. Listener
	listener := &resource.Resource{
		ID:       listenerID,
		Name:     "WebListener",
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

	// Sort resources
	graph := architecture.NewGraph(arch)
	sortedResources, err := graph.GetSortedResources()
	if err != nil {
		return fmt.Errorf("getting sorted resources: %w", err)
	}

	// Generate Terraform
	// Wire Terraform mapper registry
	mapperRegistry := tfmapper.NewRegistry()
	if err := mapperRegistry.Register(terraform.New()); err != nil {
		return fmt.Errorf("register aws terraform mapper: %w", err)
	}

	engine := tfgen.NewEngine(mapperRegistry)
	output, err := engine.Generate(ctx, arch, sortedResources)
	if err != nil {
		return fmt.Errorf("terraform engine generate: %w", err)
	}

	// Print output
	fmt.Println("Generated Terraform Files:")
	for _, f := range output.Files {
		fmt.Printf("--- %s ---\n", f.Path)
		fmt.Println(f.Content)
	}

	// Write to file for inspection
	outDir := filepath.Join("terraform_output", "scenario17")
	if err := writeTerraformOutput(outDir, output); err != nil {
		return err
	}
	fmt.Printf("\nFiles written to %s\n", outDir)

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
