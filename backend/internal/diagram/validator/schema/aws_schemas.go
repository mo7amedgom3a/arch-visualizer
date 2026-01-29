package schema

// RegisterAWSSchemas registers all AWS resource schemas
func RegisterAWSSchemas(registry *InMemorySchemaRegistry) {
	// Region schema (special - container only)
	registry.Register(&ResourceSchema{
		ResourceType: "region",
		Provider:     "aws",
		Category:     "global",
		Description:  "AWS Region - top-level container for all resources",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "Region name (e.g., us-east-1)"},
		},
		ValidChildTypes: []string{"vpc"},
	})

	// VPC schema
	registry.Register(&ResourceSchema{
		ResourceType: "vpc",
		Provider:     "aws",
		Category:     "networking",
		Description:  "Virtual Private Cloud",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "VPC name"},
			{Name: "cidr", Type: FieldTypeCIDR, Required: true, Description: "VPC CIDR block", Constraints: &FieldConstraint{CIDRVersion: strPtr("ipv4")}},
			{Name: "regionId", Type: FieldTypeString, Required: false, Description: "Region ID reference"},
			{Name: "enable_dns_hostnames", Type: FieldTypeBool, Required: false, Description: "Enable DNS hostnames"},
			{Name: "enable_dns_support", Type: FieldTypeBool, Required: false, Description: "Enable DNS support"},
			{Name: "instance_tenancy", Type: FieldTypeString, Required: false, Description: "Instance tenancy", Constraints: &FieldConstraint{Enum: []string{"default", "dedicated"}}},
		},
		ValidParentTypes: []string{"region"},
		ValidChildTypes:  []string{"subnet", "security-group", "route-table", "internet-gateway", "nat-gateway"},
	})

	// Subnet schema
	registry.Register(&ResourceSchema{
		ResourceType: "subnet",
		Provider:     "aws",
		Category:     "networking",
		Description:  "VPC Subnet",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "Subnet name"},
			{Name: "cidr", Type: FieldTypeCIDR, Required: true, Description: "Subnet CIDR block", Constraints: &FieldConstraint{CIDRVersion: strPtr("ipv4")}},
			{Name: "availabilityZoneId", Type: FieldTypeString, Required: true, Description: "Availability Zone (e.g., us-east-1a)"},
			{Name: "vpcId", Type: FieldTypeString, Required: false, Description: "VPC ID reference"},
			{Name: "map_public_ip_on_launch", Type: FieldTypeBool, Required: false, Description: "Auto-assign public IP"},
		},
		ValidParentTypes: []string{"vpc"},
		ValidChildTypes:  []string{"ec2", "nat-gateway", "load-balancer"},
	})

	// EC2 Instance schema
	registry.Register(&ResourceSchema{
		ResourceType: "ec2",
		Provider:     "aws",
		Category:     "compute",
		Description:  "EC2 Instance",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "Instance name"},
			{Name: "ami", Type: FieldTypeString, Required: true, Description: "AMI ID", Constraints: &FieldConstraint{Prefix: strPtr("ami-"), MinLength: intPtr(12), MaxLength: intPtr(21)}},
			{Name: "instanceType", Type: FieldTypeString, Required: true, Description: "Instance type (e.g., t3.micro)"},
			{Name: "subnetId", Type: FieldTypeString, Required: false, Description: "Subnet ID reference"},
			{Name: "securityGroupIds", Type: FieldTypeArray, Required: false, Description: "Security group IDs", ItemType: fieldTypePtr(FieldTypeString)},
			{Name: "keyName", Type: FieldTypeString, Required: false, Description: "SSH key pair name"},
			{Name: "iamInstanceProfile", Type: FieldTypeString, Required: false, Description: "IAM instance profile"},
			{Name: "userData", Type: FieldTypeString, Required: false, Description: "User data script", Constraints: &FieldConstraint{MaxLength: intPtr(12288)}},
		},
		ValidParentTypes: []string{"subnet"},
		ValidChildTypes:  []string{},
	})

	// Security Group schema
	registry.Register(&ResourceSchema{
		ResourceType: "security-group",
		Provider:     "aws",
		Category:     "networking",
		Description:  "Security Group",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "Security group name"},
			{Name: "description", Type: FieldTypeString, Required: false, Description: "Security group description", Constraints: &FieldConstraint{MaxLength: intPtr(255)}},
			{Name: "vpcId", Type: FieldTypeString, Required: false, Description: "VPC ID reference"},
			{Name: "ingressRules", Type: FieldTypeArray, Required: false, Description: "Inbound rules"},
			{Name: "egressRules", Type: FieldTypeArray, Required: false, Description: "Outbound rules"},
		},
		ValidParentTypes: []string{"vpc"},
		ValidChildTypes:  []string{},
	})

	// Route Table schema
	registry.Register(&ResourceSchema{
		ResourceType: "route-table",
		Provider:     "aws",
		Category:     "networking",
		Description:  "Route Table",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: false, Description: "Route table name"},
			{Name: "vpcId", Type: FieldTypeString, Required: false, Description: "VPC ID reference"},
			{Name: "isMain", Type: FieldTypeBool, Required: false, Description: "Is main route table"},
			{Name: "routes", Type: FieldTypeArray, Required: false, Description: "Route entries"},
		},
		ValidParentTypes: []string{"vpc"},
		ValidChildTypes:  []string{},
	})

	// Internet Gateway schema
	registry.Register(&ResourceSchema{
		ResourceType: "internet-gateway",
		Provider:     "aws",
		Category:     "networking",
		Description:  "Internet Gateway",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: false, Description: "Internet gateway name"},
			{Name: "vpcId", Type: FieldTypeString, Required: false, Description: "VPC ID reference"},
		},
		ValidParentTypes: []string{"vpc"},
		ValidChildTypes:  []string{},
	})

	// NAT Gateway schema
	registry.Register(&ResourceSchema{
		ResourceType: "nat-gateway",
		Provider:     "aws",
		Category:     "networking",
		Description:  "NAT Gateway",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: false, Description: "NAT gateway name"},
			{Name: "subnetId", Type: FieldTypeString, Required: true, Description: "Subnet ID (must be public subnet)"},
			{Name: "allocationId", Type: FieldTypeString, Required: false, Description: "Elastic IP allocation ID"},
			{Name: "connectivityType", Type: FieldTypeString, Required: false, Description: "Connectivity type", Constraints: &FieldConstraint{Enum: []string{"public", "private"}}},
		},
		ValidParentTypes: []string{"subnet", "vpc"},
		ValidChildTypes:  []string{},
	})

	// Elastic IP schema
	registry.Register(&ResourceSchema{
		ResourceType: "elastic-ip",
		Provider:     "aws",
		Category:     "networking",
		Description:  "Elastic IP Address",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: false, Description: "Elastic IP name"},
			{Name: "domain", Type: FieldTypeString, Required: false, Description: "Domain (vpc or standard)", Constraints: &FieldConstraint{Enum: []string{"vpc", "standard"}}},
			{Name: "instanceId", Type: FieldTypeString, Required: false, Description: "Instance ID to associate"},
		},
		ValidParentTypes: []string{"vpc", "region"},
		ValidChildTypes:  []string{},
	})

	// Lambda schema
	registry.Register(&ResourceSchema{
		ResourceType: "lambda",
		Provider:     "aws",
		Category:     "compute",
		Description:  "Lambda Function",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "Function name"},
			{Name: "runtime", Type: FieldTypeString, Required: true, Description: "Runtime (e.g., nodejs18.x, python3.11)"},
			{Name: "handler", Type: FieldTypeString, Required: true, Description: "Handler function"},
			{Name: "memorySize", Type: FieldTypeInt, Required: false, Description: "Memory in MB", Constraints: &FieldConstraint{MinValue: floatPtr(128), MaxValue: floatPtr(10240)}},
			{Name: "timeout", Type: FieldTypeInt, Required: false, Description: "Timeout in seconds", Constraints: &FieldConstraint{MinValue: floatPtr(1), MaxValue: floatPtr(900)}},
			{Name: "role", Type: FieldTypeString, Required: true, Description: "IAM role ARN"},
		},
		ValidParentTypes: []string{"vpc", "region"},
		ValidChildTypes:  []string{},
	})

	// S3 Bucket schema
	registry.Register(&ResourceSchema{
		ResourceType: "s3",
		Provider:     "aws",
		Category:     "storage",
		Description:  "S3 Bucket",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "Bucket name (globally unique)", Constraints: &FieldConstraint{MinLength: intPtr(3), MaxLength: intPtr(63)}},
			{Name: "acl", Type: FieldTypeString, Required: false, Description: "Access control list", Constraints: &FieldConstraint{Enum: []string{"private", "public-read", "public-read-write", "authenticated-read"}}},
			{Name: "versioning", Type: FieldTypeBool, Required: false, Description: "Enable versioning"},
			{Name: "encryption", Type: FieldTypeString, Required: false, Description: "Server-side encryption", Constraints: &FieldConstraint{Enum: []string{"AES256", "aws:kms"}}},
		},
		ValidParentTypes: []string{"region"},
		ValidChildTypes:  []string{},
	})

	// RDS schema
	registry.Register(&ResourceSchema{
		ResourceType: "rds",
		Provider:     "aws",
		Category:     "database",
		Description:  "RDS Database Instance",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "DB instance identifier"},
			{Name: "engine", Type: FieldTypeString, Required: true, Description: "Database engine", Constraints: &FieldConstraint{Enum: []string{"mysql", "postgres", "mariadb", "oracle-se2", "sqlserver-se"}}},
			{Name: "engineVersion", Type: FieldTypeString, Required: false, Description: "Engine version"},
			{Name: "instanceClass", Type: FieldTypeString, Required: true, Description: "DB instance class (e.g., db.t3.micro)"},
			{Name: "allocatedStorage", Type: FieldTypeInt, Required: true, Description: "Storage in GB", Constraints: &FieldConstraint{MinValue: floatPtr(20), MaxValue: floatPtr(65536)}},
			{Name: "username", Type: FieldTypeString, Required: true, Description: "Master username"},
			{Name: "password", Type: FieldTypeString, Required: false, Description: "Master password"},
			{Name: "multiAz", Type: FieldTypeBool, Required: false, Description: "Multi-AZ deployment"},
		},
		ValidParentTypes: []string{"subnet", "vpc"},
		ValidChildTypes:  []string{},
	})

	// DynamoDB schema
	registry.Register(&ResourceSchema{
		ResourceType: "dynamodb",
		Provider:     "aws",
		Category:     "database",
		Description:  "DynamoDB Table",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "Table name"},
			{Name: "hashKey", Type: FieldTypeString, Required: true, Description: "Partition key name"},
			{Name: "hashKeyType", Type: FieldTypeString, Required: false, Description: "Partition key type", Constraints: &FieldConstraint{Enum: []string{"S", "N", "B"}}},
			{Name: "rangeKey", Type: FieldTypeString, Required: false, Description: "Sort key name"},
			{Name: "rangeKeyType", Type: FieldTypeString, Required: false, Description: "Sort key type", Constraints: &FieldConstraint{Enum: []string{"S", "N", "B"}}},
			{Name: "billingMode", Type: FieldTypeString, Required: false, Description: "Billing mode", Constraints: &FieldConstraint{Enum: []string{"PROVISIONED", "PAY_PER_REQUEST"}}},
			{Name: "readCapacity", Type: FieldTypeInt, Required: false, Description: "Read capacity units"},
			{Name: "writeCapacity", Type: FieldTypeInt, Required: false, Description: "Write capacity units"},
		},
		ValidParentTypes: []string{"region"},
		ValidChildTypes:  []string{},
	})

	// Load Balancer schema
	registry.Register(&ResourceSchema{
		ResourceType: "load-balancer",
		Provider:     "aws",
		Category:     "compute",
		Description:  "Application/Network Load Balancer",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: true, Description: "Load balancer name"},
			{Name: "type", Type: FieldTypeString, Required: false, Description: "Load balancer type", Constraints: &FieldConstraint{Enum: []string{"application", "network", "gateway"}}},
			{Name: "internal", Type: FieldTypeBool, Required: false, Description: "Internal load balancer"},
			{Name: "subnets", Type: FieldTypeArray, Required: false, Description: "Subnet IDs", ItemType: fieldTypePtr(FieldTypeString)},
			{Name: "securityGroups", Type: FieldTypeArray, Required: false, Description: "Security group IDs (ALB only)", ItemType: fieldTypePtr(FieldTypeString)},
		},
		ValidParentTypes: []string{"subnet", "vpc"},
		ValidChildTypes:  []string{},
	})

	// Auto Scaling Group schema
	registry.Register(&ResourceSchema{
		ResourceType: "auto-scaling-group",
		Provider:     "aws",
		Category:     "compute",
		Description:  "Auto Scaling Group",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: false, Description: "Auto scaling group name"},
			{Name: "minSize", Type: FieldTypeInt, Required: true, Description: "Minimum size", Constraints: &FieldConstraint{MinValue: floatPtr(0)}},
			{Name: "maxSize", Type: FieldTypeInt, Required: true, Description: "Maximum size", Constraints: &FieldConstraint{MinValue: floatPtr(1)}},
			{Name: "desiredCapacity", Type: FieldTypeInt, Required: false, Description: "Desired capacity"},
			{Name: "launchTemplateId", Type: FieldTypeString, Required: false, Description: "Launch template ID"},
			{Name: "subnets", Type: FieldTypeArray, Required: true, Description: "Subnet IDs", ItemType: fieldTypePtr(FieldTypeString)},
			{Name: "healthCheckType", Type: FieldTypeString, Required: false, Description: "Health check type", Constraints: &FieldConstraint{Enum: []string{"EC2", "ELB"}}},
		},
		ValidParentTypes: []string{"vpc"},
		ValidChildTypes:  []string{},
	})

	// EBS Volume schema
	registry.Register(&ResourceSchema{
		ResourceType: "ebs",
		Provider:     "aws",
		Category:     "storage",
		Description:  "EBS Volume",
		Fields: []FieldSpec{
			{Name: "name", Type: FieldTypeString, Required: false, Description: "Volume name"},
			{Name: "availabilityZone", Type: FieldTypeString, Required: true, Description: "Availability zone"},
			{Name: "size", Type: FieldTypeInt, Required: true, Description: "Size in GB", Constraints: &FieldConstraint{MinValue: floatPtr(1), MaxValue: floatPtr(16384)}},
			{Name: "volumeType", Type: FieldTypeString, Required: false, Description: "Volume type", Constraints: &FieldConstraint{Enum: []string{"gp2", "gp3", "io1", "io2", "sc1", "st1", "standard"}}},
			{Name: "iops", Type: FieldTypeInt, Required: false, Description: "IOPS (io1/io2/gp3)"},
			{Name: "throughput", Type: FieldTypeInt, Required: false, Description: "Throughput in MB/s (gp3 only)"},
			{Name: "encrypted", Type: FieldTypeBool, Required: false, Description: "Enable encryption"},
		},
		ValidParentTypes: []string{"subnet", "region"},
		ValidChildTypes:  []string{},
	})
}

// Helper functions for creating pointers
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func fieldTypePtr(ft FieldType) *FieldType {
	return &ft
}

// init registers AWS schemas on package load
func init() {
	RegisterAWSSchemas(DefaultRegistry)
}
