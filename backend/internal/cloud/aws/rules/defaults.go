package rules

// DefaultNetworkingRules returns the default AWS networking rules
// These are code-defined defaults that can be overridden by database constraints
func DefaultNetworkingRules() []ConstraintRecord {
	return []ConstraintRecord{
		// VPC Rules
		// VPC has no parent requirement (it's a top-level resource)
		// VPC requires region
		{ResourceType: "VPC", ConstraintType: "requires_region", ConstraintValue: "true"},
		// VPC can have many subnets (AWS limit is 200 per VPC, but we'll use a reasonable default)
		{ResourceType: "VPC", ConstraintType: "max_children", ConstraintValue: "200"},
		// VPC can have dependencies on InternetGateway, NATGateway
		{ResourceType: "VPC", ConstraintType: "allowed_dependencies", ConstraintValue: "InternetGateway,NATGateway"},

		// Subnet Rules
		// Subnet requires VPC as parent
		{ResourceType: "Subnet", ConstraintType: "requires_parent", ConstraintValue: "VPC"},
		// Subnet can only have VPC as parent
		{ResourceType: "Subnet", ConstraintType: "allowed_parent", ConstraintValue: "VPC"},
		// Subnet requires region (inherited from VPC, but explicit)
		{ResourceType: "Subnet", ConstraintType: "requires_region", ConstraintValue: "true"},
		// Subnet can have dependencies on RouteTable, NATGateway
		{ResourceType: "Subnet", ConstraintType: "allowed_dependencies", ConstraintValue: "RouteTable,NATGateway"},
		// Subnet cannot depend on itself or VPC (circular)
		{ResourceType: "Subnet", ConstraintType: "forbidden_dependencies", ConstraintValue: "Subnet,VPC"},

		// InternetGateway Rules
		// InternetGateway requires VPC as parent (attachment)
		{ResourceType: "InternetGateway", ConstraintType: "requires_parent", ConstraintValue: "VPC"},
		// InternetGateway can only have VPC as parent
		{ResourceType: "InternetGateway", ConstraintType: "allowed_parent", ConstraintValue: "VPC"},
		// InternetGateway requires region
		{ResourceType: "InternetGateway", ConstraintType: "requires_region", ConstraintValue: "true"},
		// InternetGateway can be used by RouteTable
		{ResourceType: "InternetGateway", ConstraintType: "allowed_dependencies", ConstraintValue: "RouteTable"},
		// InternetGateway cannot depend on Subnet or NATGateway
		{ResourceType: "InternetGateway", ConstraintType: "forbidden_dependencies", ConstraintValue: "Subnet,NATGateway"},

		// RouteTable Rules
		// RouteTable requires VPC as parent
		{ResourceType: "RouteTable", ConstraintType: "requires_parent", ConstraintValue: "VPC"},
		// RouteTable can only have VPC as parent
		{ResourceType: "RouteTable", ConstraintType: "allowed_parent", ConstraintValue: "VPC"},
		// RouteTable requires region
		{ResourceType: "RouteTable", ConstraintType: "requires_region", ConstraintValue: "true"},
		// RouteTable can depend on InternetGateway, NATGateway, Subnet
		{ResourceType: "RouteTable", ConstraintType: "allowed_dependencies", ConstraintValue: "InternetGateway,NATGateway,Subnet"},
		// RouteTable cannot depend on SecurityGroup
		{ResourceType: "RouteTable", ConstraintType: "forbidden_dependencies", ConstraintValue: "SecurityGroup"},

		// SecurityGroup Rules
		// SecurityGroup requires VPC as parent
		{ResourceType: "SecurityGroup", ConstraintType: "requires_parent", ConstraintValue: "VPC"},
		// SecurityGroup can only have VPC as parent
		{ResourceType: "SecurityGroup", ConstraintType: "allowed_parent", ConstraintValue: "VPC"},
		// SecurityGroup requires region
		{ResourceType: "SecurityGroup", ConstraintType: "requires_region", ConstraintValue: "true"},
		// SecurityGroup can depend on other SecurityGroups (for source group rules)
		{ResourceType: "SecurityGroup", ConstraintType: "allowed_dependencies", ConstraintValue: "SecurityGroup"},
		// SecurityGroup cannot depend on networking infrastructure
		{ResourceType: "SecurityGroup", ConstraintType: "forbidden_dependencies", ConstraintValue: "VPC,Subnet,RouteTable,InternetGateway,NATGateway"},

		// NATGateway Rules
		// NATGateway requires Subnet as parent
		{ResourceType: "NATGateway", ConstraintType: "requires_parent", ConstraintValue: "Subnet"},
		// NATGateway can only have Subnet as parent
		{ResourceType: "NATGateway", ConstraintType: "allowed_parent", ConstraintValue: "Subnet"},
		// NATGateway requires region
		{ResourceType: "NATGateway", ConstraintType: "requires_region", ConstraintValue: "true"},
		// NATGateway can be used by RouteTable
		{ResourceType: "NATGateway", ConstraintType: "allowed_dependencies", ConstraintValue: "RouteTable"},
		// NATGateway cannot depend on InternetGateway or VPC directly
		{ResourceType: "NATGateway", ConstraintType: "forbidden_dependencies", ConstraintValue: "InternetGateway,VPC"},
	}
}

// DefaultComputeRules returns the default AWS compute rules
func DefaultComputeRules() []ConstraintRecord {
	return []ConstraintRecord{
		// EC2 Rules
		// EC2 requires Subnet as parent
		{ResourceType: "EC2", ConstraintType: "requires_parent", ConstraintValue: "Subnet"},
		// EC2 can only have Subnet as parent
		{ResourceType: "EC2", ConstraintType: "allowed_parent", ConstraintValue: "Subnet"},
		// EC2 requires region
		{ResourceType: "EC2", ConstraintType: "requires_region", ConstraintValue: "true"},
		// EC2 can have dependencies on EBS, ENI, EIP, SecurityGroup, IAMRole
		{ResourceType: "EC2", ConstraintType: "allowed_dependencies", ConstraintValue: "EBS,NetworkInterface,ElasticIP,SecurityGroup,IAMRole"},
		// EC2 cannot depend on VPC directly
		{ResourceType: "EC2", ConstraintType: "forbidden_dependencies", ConstraintValue: "VPC"},

		// Lambda Rules
		// Lambda is regional, no strict parent requirement in our model (can be standalone or in VPC)
		// But for now, let's say it's regional (parent: nil or Region via visualizer logic?)
		// In our graph model, regional resources often have no parent or Region as parent.
		// Let's enforce Region requirement.
		{ResourceType: "Lambda", ConstraintType: "requires_region", ConstraintValue: "true"},
		// Lambda can depend on IAMRole, S3, DynamoDB, etc.
		{ResourceType: "Lambda", ConstraintType: "allowed_dependencies", ConstraintValue: "IAMRole,S3,DynamoDB,SQS,SNS"},
	}
}

// DefaultStorageRules returns the default AWS storage rules
func DefaultStorageRules() []ConstraintRecord {
	return []ConstraintRecord{
		// S3 Rules
		// S3 is global/regional mixed, but bucket exists in a region.
		{ResourceType: "S3", ConstraintType: "requires_region", ConstraintValue: "true"},
		// S3 often has no parent in visualizer (top level)

		// EBS Rules
		// EBS needs to be attached to EC2 or just exist in AZ?
		// Usually visualized attached to EC2 or independent.
		{ResourceType: "EBS", ConstraintType: "requires_region", ConstraintValue: "true"},
		// EBS allowed dependency: Snapshot?
		// Allowed parent? Usually none if independent, or EC2 if we model containment (but EBS isn't contained in EC2).
		// Let's leave parent constraints loose for EBS for now.
	}
}

// DefaultDatabaseRules returns the default AWS database rules
func DefaultDatabaseRules() []ConstraintRecord {
	return []ConstraintRecord{
		// RDS Rules
		// RDS requires Subnet as parent (must run in a VPC)
		{ResourceType: "RDS", ConstraintType: "requires_parent", ConstraintValue: "Subnet"},
		// RDS requires region
		{ResourceType: "RDS", ConstraintType: "requires_region", ConstraintValue: "true"},
		// RDS can depend on SecurityGroup, S3 (import/export), IAMRole, KMS
		{ResourceType: "RDS", ConstraintType: "allowed_dependencies", ConstraintValue: "SecurityGroup,S3,IAMRole,KMSKey"},
		// RDS cannot depend on VPC directly (must go through subnet)
		{ResourceType: "RDS", ConstraintType: "forbidden_dependencies", ConstraintValue: "VPC"},

		// DynamoDB Rules
		// DynamoDB is regional, no parent requirement usually
		{ResourceType: "DynamoDB", ConstraintType: "requires_region", ConstraintValue: "true"},
		// DynamoDB can depend on IAMRole, KMS
		{ResourceType: "DynamoDB", ConstraintType: "allowed_dependencies", ConstraintValue: "IAMRole,KMSKey"},
	}
}
