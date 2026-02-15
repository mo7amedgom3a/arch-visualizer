package inventory

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// GetAWSResourceClassifications returns all AWS resource classifications
func GetAWSResourceClassifications() []ResourceClassification {
	return []ResourceClassification{
		// Networking Resources
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "VPC",
			IRType:       "vpc",
			Aliases:      []string{"vpc"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "Subnet",
			IRType:       "subnet",
			Aliases:      []string{"subnet"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "AvailabilityZone",
			IRType:       "availability-zone",
			Aliases:      []string{"availability-zone", "az"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "RouteTable",
			IRType:       "route-table",
			Aliases:      []string{"route-table", "route_table"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "SecurityGroup",
			IRType:       "security-group",
			Aliases:      []string{"security-group", "security_group", "sg"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "InternetGateway",
			IRType:       "internet-gateway",
			Aliases:      []string{"internet-gateway", "internet_gateway", "igw"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "NATGateway",
			IRType:       "nat-gateway",
			Aliases:      []string{"nat-gateway", "nat_gateway", "nat"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "ElasticIP",
			IRType:       "elastic-ip",
			Aliases:      []string{"elastic-ip", "elastic_ip", "eip"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "VPCEndpoint",
			IRType:       "vpc-endpoint",
			Aliases:      []string{"vpc-endpoint", "vpc_endpoint", "vpce"},
		},

		// Compute Resources
		{
			Category:     resource.CategoryCompute,
			ResourceName: "EC2",
			IRType:       "ec2",
			Aliases:      []string{"ec2", "instance", "ec2-instance"},
		},
		{
			Category:     resource.CategoryCompute,
			ResourceName: "Lambda",
			IRType:       "lambda",
			Aliases:      []string{"lambda", "lambda-function"},
		},
		{
			Category:     resource.CategoryCompute,
			ResourceName: "LoadBalancer",
			IRType:       "load-balancer",
			Aliases:      []string{"load-balancer", "load_balancer", "elb", "alb", "nlb"},
		},
		{
			Category:     resource.CategoryCompute,
			ResourceName: "Listener",
			IRType:       "listener",
			Aliases:      []string{"listener", "alb-listener", "lb-listener"},
		},
		{
			Category:     resource.CategoryCompute,
			ResourceName: "TargetGroup",
			IRType:       "target-group",
			Aliases:      []string{"target-group", "tg", "alb-target-group"},
		},
		{
			Category:     resource.CategoryCompute,
			ResourceName: "AutoScalingGroup",
			IRType:       "autoscaling-group",
			Aliases:      []string{"autoscaling-group", "auto-scaling-group", "auto_scaling_group", "asg"},
		},
		{
			Category:     resource.CategoryCompute,
			ResourceName: "LaunchTemplate",
			IRType:       "launch-template",
			Aliases:      []string{"launch-template", "launch_template", "lt"},
		},

		// Storage Resources
		{
			Category:     resource.CategoryStorage,
			ResourceName: "S3",
			IRType:       "s3",
			Aliases:      []string{"s3", "s3-bucket", "bucket"},
		},
		{
			Category:     resource.CategoryStorage,
			ResourceName: "EBS",
			IRType:       "ebs",
			Aliases:      []string{"ebs", "ebs-volume", "volume"},
		},

		// Database Resources
		{
			Category:     resource.CategoryDatabase,
			ResourceName: "RDS",
			IRType:       "rds",
			Aliases:      []string{"rds", "rds-instance"},
		},
		{
			Category:     resource.CategoryDatabase,
			ResourceName: "DynamoDB",
			IRType:       "dynamodb",
			Aliases:      []string{"dynamodb", "dynamo-db"},
		},

		// IAM Resources
		{
			Category:     resource.CategoryIAM,
			ResourceName: "IAMPolicy",
			IRType:       "iam-policy",
			Aliases:      []string{"iam-policy", "aws_iam_policy"},
		},
		{
			Category:     resource.CategoryIAM,
			ResourceName: "IAMUser",
			IRType:       "iam-user",
			Aliases:      []string{"iam-user", "aws_iam_user", "user"},
		},
		{
			Category:     resource.CategoryIAM,
			ResourceName: "IAMRole",
			IRType:       "iam-role",
			Aliases:      []string{"iam-role", "aws_iam_role", "role"},
		},
		{
			Category:     resource.CategoryIAM,
			ResourceName: "IAMRolePolicyAttachment",
			IRType:       "iam-role-policy-attachment",
			Aliases:      []string{"iam-role-policy-attachment", "aws_iam_role_policy_attachment"},
		},

		{
			Category:     resource.CategoryIAM,
			ResourceName: "IAMInstanceProfile",
			IRType:       "iam-instance-profile",
			Aliases:      []string{"iam-instance-profile", "aws_iam_instance_profile", "instance-profile"},
		},

		// Container Resources (ECS)
		{
			Category:     resource.CategoryContainers,
			ResourceName: "ECRRepository",
			IRType:       "ecr-repository",
			Aliases:      []string{"ecr-repository", "ecr_repository", "ecr"},
		},
		{
			Category:     resource.CategoryContainers,
			ResourceName: "ECSCluster",
			IRType:       "ecs-cluster",
			Aliases:      []string{"ecs-cluster", "ecs_cluster"},
		},
		{
			Category:     resource.CategoryContainers,
			ResourceName: "ECSTaskDefinition",
			IRType:       "ecs-task-definition",
			Aliases:      []string{"ecs-task-definition", "task-definition", "ecs_task_definition"},
		},
		{
			Category:     resource.CategoryContainers,
			ResourceName: "ECSService",
			IRType:       "ecs-service",
			Aliases:      []string{"ecs-service", "ecs_service"},
		},
		{
			Category:     resource.CategoryContainers,
			ResourceName: "ECSCapacityProvider",
			IRType:       "ecs-capacity-provider",
			Aliases:      []string{"ecs-capacity-provider", "capacity-provider", "ecs_capacity_provider"},
		},
		{
			Category:     resource.CategoryContainers,
			ResourceName: "ECSClusterCapacityProviders",
			IRType:       "ecs-cluster-capacity-providers",
			Aliases:      []string{"ecs-cluster-capacity-providers", "ecs_cluster_capacity_providers"},
		},
	}
}
