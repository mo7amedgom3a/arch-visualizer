package architecture

import (
	"testing"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory" // Initialize inventory
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

func TestAWSArchitectureGenerator_Provider(t *testing.T) {
	generator := NewAWSArchitectureGenerator()
	if generator.Provider() != resource.AWS {
		t.Errorf("Expected provider to be AWS, got %s", generator.Provider())
	}
}

func TestAWSArchitectureGenerator_Generate_SimpleArchitecture(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	// Create a simple diagram graph with region and VPC
	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add region node
	diagramGraph.Nodes["region-1"] = &graph.Node{
		ID:           "region-1",
		Type:         "containerNode",
		ResourceType: "region",
		Label:        "US East 1",
		Config:       map[string]interface{}{"name": "us-east-1"},
		PositionX:    100,
		PositionY:    100,
	}

	// Add VPC node
	diagramGraph.Nodes["vpc-1"] = &graph.Node{
		ID:           "vpc-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "Main VPC",
		Config: map[string]interface{}{
			"name":      "main-vpc",
			"cidrBlock": "10.0.0.0/16",
		},
		ParentID:  stringPtr("region-1"),
		PositionX: 200,
		PositionY: 200,
	}

	// Generate architecture
	arch, err := generator.Generate(diagramGraph)
	if err != nil {
		t.Fatalf("Failed to generate architecture: %v", err)
	}

	// Verify architecture
	if arch == nil {
		t.Fatal("Expected non-nil architecture")
	}

	if arch.Provider != resource.AWS {
		t.Errorf("Expected provider to be AWS, got %s", arch.Provider)
	}

	if arch.Region != "us-east-1" {
		t.Errorf("Expected region to be 'us-east-1', got '%s'", arch.Region)
	}

	if len(arch.Resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(arch.Resources))
	}

	// Verify VPC resource
	vpc := arch.Resources[0]
	if vpc.ID != "vpc-1" {
		t.Errorf("Expected resource ID to be 'vpc-1', got '%s'", vpc.ID)
	}

	if vpc.Name != "main-vpc" {
		t.Errorf("Expected resource name to be 'main-vpc', got '%s'", vpc.Name)
	}

	if vpc.Type.Name != "VPC" {
		t.Errorf("Expected resource type name to be 'VPC', got '%s'", vpc.Type.Name)
	}

	if vpc.Type.Category != string(resource.CategoryNetworking) {
		t.Errorf("Expected resource category to be 'Networking', got '%s'", vpc.Type.Category)
	}

	if vpc.Region != "us-east-1" {
		t.Errorf("Expected resource region to be 'us-east-1', got '%s'", vpc.Region)
	}

	if vpc.Metadata["name"] != "main-vpc" {
		t.Error("Expected metadata to contain 'name' field")
	}
}

func TestAWSArchitectureGenerator_Generate_WithContainment(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add region
	diagramGraph.Nodes["region-1"] = &graph.Node{
		ID:           "region-1",
		Type:         "containerNode",
		ResourceType: "region",
		Label:        "US East 1",
		Config:       map[string]interface{}{"name": "us-east-1"},
	}

	// Add VPC
	diagramGraph.Nodes["vpc-1"] = &graph.Node{
		ID:           "vpc-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "Main VPC",
		Config:       map[string]interface{}{"name": "main-vpc"},
		ParentID:     stringPtr("region-1"),
	}

	// Add Subnet (child of VPC)
	diagramGraph.Nodes["subnet-1"] = &graph.Node{
		ID:           "subnet-1",
		Type:         "containerNode",
		ResourceType: "subnet",
		Label:        "Private Subnet",
		Config: map[string]interface{}{
			"name":      "private-subnet",
			"cidrBlock": "10.0.1.0/24",
		},
		ParentID: stringPtr("vpc-1"),
	}

	// Generate architecture
	arch, err := generator.Generate(diagramGraph)
	if err != nil {
		t.Fatalf("Failed to generate architecture: %v", err)
	}

	// Verify resources
	if len(arch.Resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(arch.Resources))
	}

	// Verify containment relationships
	if len(arch.Containments) != 1 {
		t.Errorf("Expected 1 containment relationship, got %d", len(arch.Containments))
	}

	children, exists := arch.Containments["vpc-1"]
	if !exists {
		t.Error("Expected VPC to have children")
	}

	if len(children) != 1 {
		t.Errorf("Expected VPC to have 1 child, got %d", len(children))
	}

	if children[0] != "subnet-1" {
		t.Errorf("Expected child to be 'subnet-1', got '%s'", children[0])
	}

	// Verify subnet has parent ID
	var subnet *resource.Resource
	for _, res := range arch.Resources {
		if res.ID == "subnet-1" {
			subnet = res
			break
		}
	}

	if subnet == nil {
		t.Fatal("Expected to find subnet resource")
	}

	if subnet.ParentID == nil {
		t.Error("Expected subnet to have parent ID")
	} else if *subnet.ParentID != "vpc-1" {
		t.Errorf("Expected subnet parent to be 'vpc-1', got '%s'", *subnet.ParentID)
	}
}

func TestAWSArchitectureGenerator_Generate_WithDependencies(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add region
	diagramGraph.Nodes["region-1"] = &graph.Node{
		ID:           "region-1",
		Type:         "containerNode",
		ResourceType: "region",
		Label:        "US East 1",
		Config:       map[string]interface{}{"name": "us-east-1"},
	}

	// Add VPC
	diagramGraph.Nodes["vpc-1"] = &graph.Node{
		ID:           "vpc-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "Main VPC",
		Config:       map[string]interface{}{"name": "main-vpc"},
		ParentID:     stringPtr("region-1"),
	}

	// Add Internet Gateway
	diagramGraph.Nodes["igw-1"] = &graph.Node{
		ID:           "igw-1",
		Type:         "resourceNode",
		ResourceType: "internet-gateway",
		Label:        "Internet Gateway",
		Config:       map[string]interface{}{"name": "main-igw"},
		ParentID:     stringPtr("vpc-1"),
	}

	// Add dependency edge: IGW depends on VPC
	diagramGraph.Edges = append(diagramGraph.Edges, &graph.Edge{
		ID:     "edge-1",
		Source: "igw-1",
		Target: "vpc-1",
		Type:   "dependency",
	})

	// Generate architecture
	arch, err := generator.Generate(diagramGraph)
	if err != nil {
		t.Fatalf("Failed to generate architecture: %v", err)
	}

	// Verify dependencies
	if len(arch.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency relationship, got %d", len(arch.Dependencies))
	}

	deps, exists := arch.Dependencies["igw-1"]
	if !exists {
		t.Error("Expected IGW to have dependencies")
	}

	if len(deps) != 1 {
		t.Errorf("Expected IGW to have 1 dependency, got %d", len(deps))
	}

	if deps[0] != "vpc-1" {
		t.Errorf("Expected IGW dependency to be 'vpc-1', got '%s'", deps[0])
	}

	// Verify IGW resource has DependsOn
	var igw *resource.Resource
	for _, res := range arch.Resources {
		if res.ID == "igw-1" {
			igw = res
			break
		}
	}

	if igw == nil {
		t.Fatal("Expected to find IGW resource")
	}

	if len(igw.DependsOn) != 1 {
		t.Errorf("Expected IGW to have 1 dependency, got %d", len(igw.DependsOn))
	}

	if igw.DependsOn[0] != "vpc-1" {
		t.Errorf("Expected IGW DependsOn to be 'vpc-1', got '%s'", igw.DependsOn[0])
	}
}

func TestAWSArchitectureGenerator_Generate_IncludesVisualOnlyNodes(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add region
	diagramGraph.Nodes["region-1"] = &graph.Node{
		ID:           "region-1",
		Type:         "containerNode",
		ResourceType: "region",
		Label:        "US East 1",
		Config:       map[string]interface{}{"name": "us-east-1"},
	}

	// Add real VPC
	diagramGraph.Nodes["vpc-1"] = &graph.Node{
		ID:           "vpc-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "Main VPC",
		Config:       map[string]interface{}{"name": "main-vpc"},
		ParentID:     stringPtr("region-1"),
		IsVisualOnly: false,
	}

	// Add visual-only node
	diagramGraph.Nodes["visual-1"] = &graph.Node{
		ID:           "visual-1",
		Type:         "resourceNode",
		ResourceType: "ec2",
		Label:        "Visual Only",
		Config:       map[string]interface{}{"name": "visual-ec2"},
		ParentID:     stringPtr("vpc-1"),
		IsVisualOnly: true,
	}

	// Generate architecture
	arch, err := generator.Generate(diagramGraph)
	if err != nil {
		t.Fatalf("Failed to generate architecture: %v", err)
	}

	// Verify visual-only node is NOT filtered out
	if len(arch.Resources) != 3 {
		t.Errorf("Expected 3 resources (VPC, visual-ec2, default-sg), got %d", len(arch.Resources))
	}

	foundVisual := false
	for _, res := range arch.Resources {
		if res.ID == "visual-1" {
			foundVisual = true
			if isVis, ok := res.Metadata["isVisualOnly"].(bool); !ok || !isVis {
				t.Errorf("Expected visual-1 metadata to have isVisualOnly=true")
			}
		}
	}
	if !foundVisual {
		t.Errorf("Expected visual-1 resource to be present")
	}
}

func TestAWSArchitectureGenerator_Generate_NoRegion(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add VPC without region
	diagramGraph.Nodes["vpc-1"] = &graph.Node{
		ID:           "vpc-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "Main VPC",
		Config:       map[string]interface{}{"name": "main-vpc"},
	}

	// Generate architecture
	arch, err := generator.Generate(diagramGraph)
	if err != nil {
		t.Fatalf("Failed to generate architecture: %v", err)
	}

	// Verify region is empty
	if arch.Region != "" {
		t.Errorf("Expected empty region, got '%s'", arch.Region)
	}

	// Verify resource still created
	if len(arch.Resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(arch.Resources))
	}
}

func TestAWSArchitectureGenerator_Generate_UnknownResourceType(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add node with unknown resource type
	diagramGraph.Nodes["unknown-1"] = &graph.Node{
		ID:           "unknown-1",
		Type:         "resourceNode",
		ResourceType: "unknown-type",
		Label:        "Unknown",
		Config:       map[string]interface{}{"name": "unknown"},
	}

	// Generate architecture - should fail
	_, err := generator.Generate(diagramGraph)
	if err == nil {
		t.Error("Expected error for unknown resource type")
	}

	if err != nil && err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestAWSArchitectureGenerator_Generate_NameFromLabel(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add node without name in config (should use label)
	diagramGraph.Nodes["vpc-1"] = &graph.Node{
		ID:           "vpc-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "My VPC",
		Config:       map[string]interface{}{"cidrBlock": "10.0.0.0/16"},
	}

	// Generate architecture
	arch, err := generator.Generate(diagramGraph)
	if err != nil {
		t.Fatalf("Failed to generate architecture: %v", err)
	}

	// Verify name comes from label
	if arch.Resources[0].Name != "My VPC" {
		t.Errorf("Expected name to be 'My VPC', got '%s'", arch.Resources[0].Name)
	}
}

func TestAWSArchitectureGenerator_Generate_NameFallback(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add node without name or label (should use fallback)
	diagramGraph.Nodes["vpc-1"] = &graph.Node{
		ID:           "vpc-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "",
		Config:       map[string]interface{}{"cidrBlock": "10.0.0.0/16"},
	}

	// Generate architecture
	arch, err := generator.Generate(diagramGraph)
	if err != nil {
		t.Fatalf("Failed to generate architecture: %v", err)
	}

	// Verify fallback name
	if arch.Resources[0].Name != "unnamed-resource" {
		t.Errorf("Expected name to be 'unnamed-resource', got '%s'", arch.Resources[0].Name)
	}
}

func TestAWSArchitectureGenerator_Generate_ComplexArchitecture(t *testing.T) {
	generator := NewAWSArchitectureGenerator()

	diagramGraph := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add region
	diagramGraph.Nodes["region-1"] = &graph.Node{
		ID:           "region-1",
		Type:         "containerNode",
		ResourceType: "region",
		Label:        "US East 1",
		Config:       map[string]interface{}{"name": "us-east-1"},
	}

	// Add VPC
	diagramGraph.Nodes["vpc-1"] = &graph.Node{
		ID:           "vpc-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "Main VPC",
		Config:       map[string]interface{}{"name": "main-vpc"},
		ParentID:     stringPtr("region-1"),
	}

	// Add Subnet
	diagramGraph.Nodes["subnet-1"] = &graph.Node{
		ID:           "subnet-1",
		Type:         "containerNode",
		ResourceType: "subnet",
		Label:        "Private Subnet",
		Config:       map[string]interface{}{"name": "private-subnet"},
		ParentID:     stringPtr("vpc-1"),
	}

	// Add Security Group
	diagramGraph.Nodes["sg-1"] = &graph.Node{
		ID:           "sg-1",
		Type:         "resourceNode",
		ResourceType: "security-group",
		Label:        "Web SG",
		Config:       map[string]interface{}{"name": "web-sg"},
		ParentID:     stringPtr("vpc-1"),
	}

	// Add EC2
	diagramGraph.Nodes["ec2-1"] = &graph.Node{
		ID:           "ec2-1",
		Type:         "resourceNode",
		ResourceType: "ec2",
		Label:        "Web Server",
		Config:       map[string]interface{}{"name": "web-server"},
		ParentID:     stringPtr("subnet-1"),
	}

	// Add dependency: EC2 depends on Security Group
	diagramGraph.Edges = append(diagramGraph.Edges, &graph.Edge{
		ID:     "edge-1",
		Source: "ec2-1",
		Target: "sg-1",
		Type:   "dependency",
	})

	// Generate architecture
	arch, err := generator.Generate(diagramGraph)
	if err != nil {
		t.Fatalf("Failed to generate architecture: %v", err)
	}

	// Verify all resources
	if len(arch.Resources) != 5 {
		t.Errorf("Expected 5 resources, got %d", len(arch.Resources))
	}

	// Verify containment relationships
	if len(arch.Containments) != 2 {
		t.Errorf("Expected 2 containment relationships, got %d", len(arch.Containments))
	}

	// Verify dependencies
	if len(arch.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency relationship, got %d", len(arch.Dependencies))
	}

	// Verify EC2 has dependency on Security Group
	var ec2 *resource.Resource
	for _, res := range arch.Resources {
		if res.ID == "ec2-1" {
			ec2 = res
			break
		}
	}

	if ec2 == nil {
		t.Fatal("Expected to find EC2 resource")
	}

	if len(ec2.DependsOn) != 1 || ec2.DependsOn[0] != "sg-1" {
		t.Errorf("Expected EC2 to depend on 'sg-1', got %v", ec2.DependsOn)
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
