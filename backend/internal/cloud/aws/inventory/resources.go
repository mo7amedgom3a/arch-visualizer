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
			ResourceName: "AutoScalingGroup",
			IRType:       "auto-scaling-group",
			Aliases:      []string{"auto-scaling-group", "auto_scaling_group", "asg"},
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
	}
}
