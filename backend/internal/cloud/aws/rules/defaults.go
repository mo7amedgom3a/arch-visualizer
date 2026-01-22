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
