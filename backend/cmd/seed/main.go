package main

import (
	"context"
	"log"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	seed "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/seed"
)

func main() {
	ctx := context.Background()

	// Seed the database
	if err := seed.SeedDatabaseWithAdapters(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Connect to database for queries
	if _, err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create database query service
	queryService, err := seed.NewDatabaseQueryService()
	if err != nil {
		log.Fatalf("Failed to create database query service: %v", err)
	}

	// Get all projects to demonstrate queries
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var projects []models.Project
	if err := db.WithContext(ctx).Find(&projects).Error; err != nil {
		log.Fatalf("Failed to query projects: %v", err)
	}

	if len(projects) == 0 {
		log.Println("No projects found in database")
		return
	}

	// Demonstrate querying and converting resources for each project
	log.Println("\n" + strings.Repeat("=", 80))
	log.Println("Demonstrating Database Queries and Domain Model Conversion")
	log.Println(strings.Repeat("=", 80))

	for _, project := range projects {
		log.Printf("\n📁 Project: %s (ID: %s)", project.Name, project.ID.String())

		// Print summary of all resources as domain models
		if err := queryService.PrintProjectResourcesSummary(ctx, project.ID); err != nil {
			log.Printf("  ⚠ Warning: failed to print summary for project %s: %v", project.Name, err)
			continue
		}

		// Demonstrate getting specific resource types
		log.Println("\n  🔍 Querying specific resource types:")

		// Get VPCs
		vpcs, err := queryService.GetVPCsAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get VPCs: %v", err)
		} else if len(vpcs) > 0 {
			log.Printf("    ✓ Found %d VPC(s)", len(vpcs))
			for _, vpc := range vpcs {
				log.Printf("      - %s (CIDR: %s, Region: %s)", vpc.Name, vpc.CIDR, vpc.Region)
			}
		}

		// Get Subnets
		subnets, err := queryService.GetSubnetsAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get Subnets: %v", err)
		} else if len(subnets) > 0 {
			log.Printf("    ✓ Found %d Subnet(s)", len(subnets))
			for _, subnet := range subnets {
				log.Printf("      - %s (CIDR: %s, VPC: %s)", subnet.Name, subnet.CIDR, subnet.VPCID)
			}
		}

		// Get EC2 Instances
		instances, err := queryService.GetInstancesAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get Instances: %v", err)
		} else if len(instances) > 0 {
			log.Printf("    ✓ Found %d EC2 Instance(s)", len(instances))
			for _, instance := range instances {
				log.Printf("      - %s (Type: %s, AMI: %s)", instance.Name, instance.InstanceType, instance.AMI)
			}
		}

		// Get Security Groups
		sgs, err := queryService.GetSecurityGroupsAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get Security Groups: %v", err)
		} else if len(sgs) > 0 {
			log.Printf("    ✓ Found %d Security Group(s)", len(sgs))
			for _, sg := range sgs {
				log.Printf("      - %s (VPC: %s)", sg.Name, sg.VPCID)
			}
		}

		// Get Load Balancers
		lbs, err := queryService.GetLoadBalancersAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get Load Balancers: %v", err)
		} else if len(lbs) > 0 {
			log.Printf("    ✓ Found %d Load Balancer(s)", len(lbs))
			for _, lb := range lbs {
				log.Printf("      - %s (Type: %s)", lb.Name, lb.Type)
			}
		}

		// Get Auto Scaling Groups
		asgs, err := queryService.GetAutoScalingGroupsAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get Auto Scaling Groups: %v", err)
		} else if len(asgs) > 0 {
			log.Printf("    ✓ Found %d Auto Scaling Group(s)", len(asgs))
			for _, asg := range asgs {
				log.Printf("      - %s (Min: %d, Max: %d, Desired: %v)",
					asg.Name, asg.MinSize, asg.MaxSize, asg.DesiredCapacity)
			}
		}

		// Get Lambda Functions
		lambdas, err := queryService.GetLambdaFunctionsAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get Lambda Functions: %v", err)
		} else if len(lambdas) > 0 {
			log.Printf("    ✓ Found %d Lambda Function(s)", len(lambdas))
			for _, lambda := range lambdas {
				arn := "N/A"
				if lambda.ARN != nil {
					arn = *lambda.ARN
				}
				log.Printf("      - %s (ARN: %s)", lambda.FunctionName, arn)
			}
		}

		// Get S3 Buckets
		buckets, err := queryService.GetS3BucketsAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get S3 Buckets: %v", err)
		} else if len(buckets) > 0 {
			log.Printf("    ✓ Found %d S3 Bucket(s)", len(buckets))
			for _, bucket := range buckets {
				log.Printf("      - %s (Region: %s)", bucket.Name, bucket.Region)
			}
		}

		// Get NAT Gateways
		nats, err := queryService.GetNATGatewaysAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get NAT Gateways: %v", err)
		} else if len(nats) > 0 {
			log.Printf("    ✓ Found %d NAT Gateway(s)", len(nats))
			for _, nat := range nats {
				log.Printf("      - %s (Subnet: %s)", nat.Name, nat.SubnetID)
			}
		}

		// Get Internet Gateways
		igws, err := queryService.GetInternetGatewaysAsDomainModels(ctx, project.ID)
		if err != nil {
			log.Printf("    ⚠ Failed to get Internet Gateways: %v", err)
		} else if len(igws) > 0 {
			log.Printf("    ✓ Found %d Internet Gateway(s)", len(igws))
			for _, igw := range igws {
				log.Printf("      - %s (VPC: %s)", igw.Name, igw.VPCID)
			}
		}
	}

	log.Println("\n" + strings.Repeat("=", 80))
	log.Println("✅ Database query and domain model conversion demonstration completed!")
	log.Println(strings.Repeat("=", 80))
}
