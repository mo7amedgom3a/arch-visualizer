package architecture

import (
	"testing"

	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// createTestResource creates a test resource with the given ID and type
func createTestResource(id, name, resourceType string) *resource.Resource {
	return &resource.Resource{
		ID:   id,
		Name: name,
		Type: resource.ResourceType{
			ID:   resourceType,
			Name: resourceType,
		},
		Provider:  resource.AWS,
		Region:    "us-east-1",
		DependsOn: []string{},
		Metadata:  make(map[string]interface{}),
	}
}

// TestTopologicalSort_SimpleDependencies tests basic dependency ordering
func TestTopologicalSort_SimpleDependencies(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS

	// Create resources: VPC -> Subnet -> EC2
	vpc := createTestResource("vpc-1", "my-vpc", "VPC")
	igw := createTestResource("igw-1", "my-igw", "InternetGateway")
	igw.ParentID = &vpc.ID
	arch.Containments["vpc-1"] = []string{"igw-1"}
	subnet := createTestResource("subnet-1", "my-subnet", "Subnet")

	ec2 := createTestResource("ec2-1", "my-ec2", "EC2")

	arch.Resources = []*resource.Resource{vpc, igw, subnet, ec2}

	// Subnet depends on VPC
	arch.Dependencies["subnet-1"] = []string{"vpc-1"}
	// IGW depends on VPC
	arch.Dependencies["igw-1"] = []string{"vpc-1"}
	// EC2 depends on Subnet
	arch.Dependencies["ec2-1"] = []string{"subnet-1"}

	graph := NewGraph(arch)
	fmt.Println("node 0", graph.architecture.Resources[0].ID)
	fmt.Println("node 1", graph.architecture.Resources[1].ID)
	fmt.Println("node 2", graph.architecture.Resources[2].ID)
	fmt.Println("node 3", graph.architecture.Resources[3].ID)
	fmt.Println("dependencies vpc-1", graph.architecture.Dependencies["vpc-1"])
	fmt.Println("dependencies subnet-1", graph.architecture.Dependencies["subnet-1"])
	fmt.Println("dependencies igw-1", graph.architecture.Dependencies["igw-1"])
	fmt.Println("dependencies ec2-1", graph.architecture.Dependencies["ec2-1"])	
	fmt.Println("containments vpc-1", graph.architecture.Containments["vpc-1"])
	fmt.Println(strings.Repeat("=", 100))
	result, err := graph.TopologicalSort()
	fmt.Println("result", result.Resources[0].ID)
	fmt.Println("result", result.Resources[1].ID)
	fmt.Println("result", result.Resources[2].ID)
	fmt.Println("result", result.Resources[3].ID)
	fmt.Println("levels", result.Levels)
	fmt.Println("has cycle", result.HasCycle)
	fmt.Println("cycle info", result.CycleInfo)
	fmt.Println(strings.Repeat("=", 100))
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if result.HasCycle {
		t.Fatal("Expected no cycle, but cycle was detected")
	}

	if len(result.Resources) != len(arch.Resources) {
		t.Fatalf("Expected %d resources, got %d", len(arch.Resources), len(result.Resources))
	}

	// Verify relative order constraints:
	// - VPC before Subnet, IGW, and EC2
	// - Subnet before EC2
	vpcIndex := -1
	subnetIndex := -1
	ec2Index := -1
	igwIndex := -1
	for i, res := range result.Resources {
		switch res.ID {
		case "vpc-1":
			vpcIndex = i
		case "subnet-1":
			subnetIndex = i
		case "ec2-1":
			ec2Index = i
		case "igw-1":
			igwIndex = i
		}
	}

	if vpcIndex == -1 || subnetIndex == -1 || ec2Index == -1 || igwIndex == -1 {
		t.Fatalf("Expected all resources except cycle to be present in sorted result, got: %#v", result.Resources)
	}

	if !(vpcIndex < subnetIndex && vpcIndex < igwIndex && vpcIndex < ec2Index) {
		t.Errorf("Expected VPC to come before subnet, igw, and ec2; got indices vpc=%d, subnet=%d, igw=%d, ec2=%d",
			vpcIndex, subnetIndex, igwIndex, ec2Index)
	}

	if !(subnetIndex < ec2Index) {
		t.Errorf("Expected subnet to come before ec2; got indices subnet=%d, ec2=%d", subnetIndex, ec2Index)
	}

	// Verify levels
	if len(result.Levels) != 3 {
		t.Errorf("Expected 3 levels, got %d", len(result.Levels))
	}

	if !(subnetIndex < ec2Index) {
		t.Errorf("Expected subnet to come before ec2; got indices subnet=%d, ec2=%d", subnetIndex, ec2Index)
	}
}

// TestTopologicalSort_ContainmentRelationships tests containment-based dependencies
func TestTopologicalSort_ContainmentRelationships(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS

	// Create resources: VPC contains Subnet, Subnet contains EC2
	vpc := createTestResource("vpc-1", "my-vpc", "VPC")
	igw := createTestResource("igw-1", "my-igw", "InternetGateway")
	igw.ParentID = &vpc.ID
	arch.Containments["vpc-1"] = []string{"igw-1"}
	subnet := createTestResource("subnet-1", "my-subnet", "Subnet")
	ec2 := createTestResource("ec2-1", "my-ec2", "EC2")

	// Set parent relationships
	parentID := "vpc-1"
	subnet.ParentID = &parentID
	subnetParentID := "subnet-1"
	ec2.ParentID = &subnetParentID

	arch.Resources = []*resource.Resource{vpc, igw, subnet, ec2}

	// Set containment relationships
	arch.Containments["vpc-1"] = []string{"subnet-1"}
	arch.Containments["subnet-1"] = []string{"ec2-1"}
	// IGW depends on VPC
	arch.Dependencies["igw-1"] = []string{"vpc-1"}

	graph := NewGraph(arch)
	result, err := graph.TopologicalSort()

	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if result.HasCycle {
		t.Fatal("Expected no cycle, but cycle was detected")
	}

	// Verify order: VPC -> IGW -> Subnet -> EC2
	if result.Resources[0].ID != "vpc-1" {
		t.Errorf("Expected first resource to be vpc-1, got %s", result.Resources[0].ID)
	}
	if result.Resources[1].ID != "igw-1" {
		t.Errorf("Expected second resource to be igw-1, got %s", result.Resources[1].ID)
	}
	if result.Resources[2].ID != "subnet-1" {
		t.Errorf("Expected second resource to be subnet-1, got %s", result.Resources[1].ID)
	}
	if result.Resources[3].ID != "ec2-1" {
		t.Errorf("Expected fourth resource to be ec2-1, got %s", result.Resources[3].ID)
	}
}

// TestTopologicalSort_MixedDependencies tests both explicit dependencies and containment
func TestTopologicalSort_MixedDependencies(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS

	// Create resources: Provider -> VPC -> Subnet -> RouteTable -> SecurityGroup -> EC2
	vpc := createTestResource("vpc-1", "my-vpc", "VPC")
	subnet := createTestResource("subnet-1", "my-subnet", "Subnet")
	routeTable := createTestResource("rt-1", "my-rt", "RouteTable")
	securityGroup := createTestResource("sg-1", "my-sg", "SecurityGroup")
	ec2 := createTestResource("ec2-1", "my-ec2", "EC2")

	// Subnet is contained in VPC
	parentID := "vpc-1"
	subnet.ParentID = &parentID
	arch.Containments["vpc-1"] = []string{"subnet-1"}

	// RouteTable depends on VPC and Subnet
	arch.Dependencies["rt-1"] = []string{"vpc-1", "subnet-1"}

	// SecurityGroup depends on VPC
	arch.Dependencies["sg-1"] = []string{"vpc-1"}

	// EC2 depends on Subnet and SecurityGroup
	arch.Dependencies["ec2-1"] = []string{"subnet-1", "sg-1"}

	arch.Resources = []*resource.Resource{vpc, subnet, routeTable, securityGroup, ec2}

	graph := NewGraph(arch)
	result, err := graph.TopologicalSort()

	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if result.HasCycle {
		t.Fatal("Expected no cycle, but cycle was detected")
	}

	if len(result.Resources) != len(arch.Resources) {
		t.Fatalf("Expected %d resources, got %d", len(arch.Resources), len(result.Resources))
	}

	// Verify VPC comes first
	if result.Resources[0].ID != "vpc-1" {
		t.Errorf("Expected first resource to be vpc-1, got %s", result.Resources[0].ID)
	}

	// Verify Subnet comes before RouteTable, SecurityGroup, and EC2
	vpcIndex := -1
	subnetIndex := -1
	rtIndex := -1
	sgIndex := -1
	ec2Index := -1

	for i, res := range result.Resources {
		switch res.ID {
		case "vpc-1":
			vpcIndex = i
		case "subnet-1":
			subnetIndex = i
		case "rt-1":
			rtIndex = i
		case "sg-1":
			sgIndex = i
		case "ec2-1":
			ec2Index = i
		}
	}

	if subnetIndex <= vpcIndex {
		t.Errorf("Expected subnet to come after vpc, but subnet=%d, vpc=%d", subnetIndex, vpcIndex)
	}
	if rtIndex <= subnetIndex {
		t.Errorf("Expected route table to come after subnet, but rt=%d, subnet=%d", rtIndex, subnetIndex)
	}
	if sgIndex <= vpcIndex {
		t.Errorf("Expected security group to come after vpc, but sg=%d, vpc=%d", sgIndex, vpcIndex)
	}
	if ec2Index <= subnetIndex {
		t.Errorf("Expected ec2 to come after subnet, but ec2=%d, subnet=%d", ec2Index, subnetIndex)
	}
	if ec2Index <= sgIndex {
		t.Errorf("Expected ec2 to come after security group, but ec2=%d, sg=%d", ec2Index, sgIndex)
	}
}

// TestTopologicalSort_CircularDependency tests cycle detection
func TestTopologicalSort_CircularDependency(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS

	// Create circular dependency: A -> B -> C -> A
	resourceA := createTestResource("a", "resource-a", "VPC")
	resourceB := createTestResource("b", "resource-b", "Subnet")
	resourceC := createTestResource("c", "resource-c", "EC2")

	arch.Resources = []*resource.Resource{resourceA, resourceB, resourceC}

	// Create cycle: A -> B -> C -> A
	arch.Dependencies["a"] = []string{"c"}
	arch.Dependencies["b"] = []string{"a"}
	arch.Dependencies["c"] = []string{"b"}

	graph := NewGraph(arch)
	result, err := graph.TopologicalSort()

	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if !result.HasCycle {
		t.Fatal("Expected cycle to be detected, but none was found")
	}

	if len(result.CycleInfo) == 0 {
		t.Error("Expected cycle info to contain resource IDs, but it's empty")
	}

	// Verify all resources are in cycle info
	cycleSet := make(map[string]bool)
	for _, id := range result.CycleInfo {
		cycleSet[id] = true
	}

	if !cycleSet["a"] || !cycleSet["b"] || !cycleSet["c"] {
		t.Errorf("Expected all resources (a, b, c) in cycle info, got: %v", result.CycleInfo)
	}
}

// TestTopologicalSort_EmptyArchitecture tests empty architecture
func TestTopologicalSort_EmptyArchitecture(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS
	arch.Resources = []*resource.Resource{}

	graph := NewGraph(arch)
	result, err := graph.TopologicalSort()

	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if result.HasCycle {
		t.Error("Empty architecture should not have cycles")
	}

	if len(result.Resources) != 0 {
		t.Errorf("Expected 0 resources, got %d", len(result.Resources))
	}

	if len(result.Levels) != 0 {
		t.Errorf("Expected 0 levels, got %d", len(result.Levels))
	}
}

// TestTopologicalSort_NoDependencies tests resources with no dependencies
func TestTopologicalSort_NoDependencies(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS

	resource1 := createTestResource("r1", "resource-1", "VPC")
	resource2 := createTestResource("r2", "resource-2", "VPC")
	resource3 := createTestResource("r3", "resource-3", "VPC")

	arch.Resources = []*resource.Resource{resource1, resource2, resource3}

	graph := NewGraph(arch)
	result, err := graph.TopologicalSort()

	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if result.HasCycle {
		t.Fatal("Expected no cycle, but cycle was detected")
	}

	if len(result.Resources) != 3 {
		t.Fatalf("Expected 3 resources, got %d", len(result.Resources))
	}

	// All resources should be in the first level (no dependencies)
	if len(result.Levels) < 1 {
		t.Fatal("Expected at least 1 level")
	}

	firstLevelSize := len(result.Levels[0])
	if firstLevelSize != 3 {
		t.Errorf("Expected all 3 resources in first level, got %d", firstLevelSize)
	}
}

// TestGetSortedResources tests the convenience method
func TestGetSortedResources(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS

	vpc := createTestResource("vpc-1", "my-vpc", "VPC")
	subnet := createTestResource("subnet-1", "my-subnet", "Subnet")
	ec2 := createTestResource("ec2-1", "my-ec2", "EC2")

	arch.Resources = []*resource.Resource{vpc, subnet, ec2}
	arch.Dependencies["subnet-1"] = []string{"vpc-1"}
	arch.Dependencies["ec2-1"] = []string{"subnet-1"}

	graph := NewGraph(arch)
	sorted, err := graph.GetSortedResources()

	if err != nil {
		t.Fatalf("GetSortedResources failed: %v", err)
	}

	if len(sorted) != 3 {
		t.Fatalf("Expected 3 resources, got %d", len(sorted))
	}

	if sorted[0].ID != "vpc-1" || sorted[1].ID != "subnet-1" || sorted[2].ID != "ec2-1" {
		t.Errorf("Unexpected order: %v", sorted)
	}
}

// TestGetSortedResources_CycleDetection tests error handling for cycles
func TestGetSortedResources_CycleDetection(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS

	resourceA := createTestResource("a", "resource-a", "VPC")
	resourceB := createTestResource("b", "resource-b", "Subnet")

	arch.Resources = []*resource.Resource{resourceA, resourceB}

	// Create cycle: A -> B -> A
	arch.Dependencies["a"] = []string{"b"}
	arch.Dependencies["b"] = []string{"a"}

	graph := NewGraph(arch)
	_, err := graph.GetSortedResources()

	if err == nil {
		t.Fatal("Expected error for circular dependency, but got nil")
	}

	// Verify error message mentions cycle
	if err.Error() == "" {
		t.Error("Expected error message, but got empty string")
	}
}

// TestTopologicalSort_ComplexAWSArchitecture tests a realistic AWS architecture
func TestTopologicalSort_ComplexAWSArchitecture(t *testing.T) {
	arch := NewArchitecture()
	arch.Provider = resource.AWS

	// Create resources in typical AWS order
	vpc := createTestResource("vpc-1", "main-vpc", "VPC")
	igw := createTestResource("igw-1", "main-igw", "InternetGateway")
	subnet1 := createTestResource("subnet-1", "public-subnet", "Subnet")
	subnet2 := createTestResource("subnet-2", "private-subnet", "Subnet")
	rt1 := createTestResource("rt-1", "public-rt", "RouteTable")
	rt2 := createTestResource("rt-2", "private-rt", "RouteTable")
	sg1 := createTestResource("sg-1", "web-sg", "SecurityGroup")
	ec2 := createTestResource("ec2-1", "web-server", "EC2")

	arch.Resources = []*resource.Resource{vpc, igw, subnet1, subnet2, rt1, rt2, sg1, ec2}

	// Set containment: VPC contains subnets
	parentID := "vpc-1"
	subnet1.ParentID = &parentID
	subnet2.ParentID = &parentID
	arch.Containments["vpc-1"] = []string{"subnet-1", "subnet-2"}

	// IGW depends on VPC
	arch.Dependencies["igw-1"] = []string{"vpc-1"}

	// RouteTables depend on VPC
	arch.Dependencies["rt-1"] = []string{"vpc-1"}
	arch.Dependencies["rt-2"] = []string{"vpc-1"}

	// SecurityGroup depends on VPC
	arch.Dependencies["sg-1"] = []string{"vpc-1"}

	// EC2 depends on Subnet and SecurityGroup
	arch.Dependencies["ec2-1"] = []string{"subnet-1", "sg-1"}

	graph := NewGraph(arch)
	result, err := graph.TopologicalSort()

	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if result.HasCycle {
		t.Fatal("Expected no cycle, but cycle was detected")
	}

	if len(result.Resources) != 8 {
		t.Fatalf("Expected 8 resources, got %d", len(result.Resources))
	}

	// Verify VPC comes first
	if result.Resources[0].ID != "vpc-1" {
		t.Errorf("Expected VPC to be first, got %s", result.Resources[0].ID)
	}

	// Verify subnets come before EC2
	subnet1Index := -1
	subnet2Index := -1
	ec2Index := -1

	for i, res := range result.Resources {
		switch res.ID {
		case "subnet-1":
			subnet1Index = i
		case "subnet-2":
			subnet2Index = i
		case "ec2-1":
			ec2Index = i
		}
	}

	if subnet1Index == -1 || subnet2Index == -1 || ec2Index == -1 {
		t.Fatal("Could not find all required resources in sorted list")
	}

	if ec2Index <= subnet1Index || ec2Index <= subnet2Index {
		t.Errorf("EC2 should come after subnets: ec2=%d, subnet1=%d, subnet2=%d", ec2Index, subnet1Index, subnet2Index)
	}
}
