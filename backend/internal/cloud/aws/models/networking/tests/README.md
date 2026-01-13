# AWS Networking Models - Unit Tests

This directory contains comprehensive unit tests for AWS networking models and constraints. The tests verify the logic and validation rules for VPCs, subnets, route tables, gateways, and their relationships.

## Overview

The test suite validates:
- **VPC Region Constraints**: Ensures VPCs are confined to a single region
- **Subnet CIDR Validation**: Checks for CIDR overlaps and validates subnet CIDRs are within VPC CIDR ranges
- **Subnet AZ Constraints**: Verifies subnets are assigned to a single Availability Zone
- **Internet Gateway**: Validates IGW configuration and VPC attachments
- **NAT Gateway**: Validates NAT gateway configuration and subnet relationships
- **Route Tables**: Tests route table validation and route target configurations
- **Route Table Associations**: Validates route table associations with subnets

## Test Structure

All tests follow a **table-driven design** pattern, making them easy to extend and maintain. Each test file:
- Defines test cases in a structured format
- Provides detailed console output showing test execution
- Validates both positive and negative scenarios
- Includes descriptive test names and error messages

## Running the Tests

### Run All Tests
```bash
go test ./internal/cloud/aws/models/networking/tests/... -v
```

### Run Specific Test
```bash
go test ./internal/cloud/aws/models/networking/tests/... -v -run TestVPCRegionConstraint
```

### Run with Coverage
```bash
go test ./internal/cloud/aws/models/networking/tests/... -cover
```

## Test Files

### 1. `vpc_test.go`
**Purpose**: Basic VPC validation tests

**Test Function**: `TestVPC`

**Coverage**:
- Valid VPC creation
- Empty name validation
- Empty region validation
- Empty CIDR validation
- Invalid CIDR format validation

### 2. `vpc_region_constraint_test.go`
**Purpose**: Tests VPC region constraints and subnet region alignment

**Test Function**: `TestVPCRegionConstraint`

**Coverage**:
- VPC in `us-east-1` with subnets in `us-east-1a` and `us-east-1b` (valid)
- VPC with empty region (invalid)
- VPC with invalid region format
- Subnet AZ validation within VPC region

**Test Scenario**:
- Region: `us-east-1`
- Availability Zones: `us-east-1a`, `us-east-1b`
- VPC CIDR: `10.0.0.0/16`

### 3. `subnet_cidr_overlap_test.go`
**Purpose**: Detects CIDR block overlaps between subnets

**Test Function**: `TestSubnetCIDROverlap`

**Coverage**:
- Four non-overlapping subnets (10.0.1.0/24, 10.0.2.0/24, 10.0.3.0/24, 10.0.4.0/24)
- Overlapping CIDR blocks (same CIDR)
- Partial CIDR overlaps
- One subnet containing another
- Adjacent non-overlapping CIDRs

**Test Scenario**:
- Public subnet: `10.0.1.0/24` (us-east-1a)
- Private subnet 1: `10.0.2.0/24` (us-east-1a)
- Private subnet 2: `10.0.3.0/24` (us-east-1b)
- Private subnet 3: `10.0.4.0/24` (us-east-1a)

### 4. `subnet_cidr_in_vpc_test.go`
**Purpose**: Validates subnet CIDR blocks are within VPC CIDR range

**Test Function**: `TestSubnetCIDRInVPC`

**Coverage**:
- Valid subnet CIDRs within VPC CIDR (10.0.0.0/16)
- Invalid subnet CIDR outside VPC CIDR
- Invalid subnet CIDR with less specific mask than VPC
- Invalid subnet CIDR with same mask as VPC
- Edge cases (10.0.255.0/24 within 10.0.0.0/16)

**Test Scenario**:
- VPC CIDR: `10.0.0.0/16`
- Valid subnets: `10.0.1.0/24`, `10.0.2.0/24`, `10.0.3.0/24`, `10.0.4.0/24`
- Invalid subnets: `172.16.0.0/24`, `192.168.1.0/24`, `10.0.0.0/8`, `10.0.0.0/16`

### 5. `subnet_az_constraint_test.go`
**Purpose**: Validates subnet Availability Zone constraints

**Test Function**: `TestSubnetAZConstraint`

**Coverage**:
- Valid subnets in `us-east-1a` and `us-east-1b`
- Invalid subnet with empty AZ
- Invalid AZ format (validation checks)
- AZ format validation (region-az pattern)

**Key Constraint**: Each subnet must be assigned to a single Availability Zone and cannot span multiple AZs.

### 6. `internet_gateway_test.go`
**Purpose**: Tests Internet Gateway validation and VPC relationships

**Test Function**: `TestInternetGateway`

**Coverage**:
- Valid IGW attached to VPC
- Invalid IGW with missing name
- Invalid IGW with missing VPC ID
- Valid IGW with multiple tags
- IGW-VPC relationship validation

### 7. `nat_gateway_test.go`
**Purpose**: Tests NAT Gateway validation and subnet relationships

**Test Function**: `TestNATGateway`

**Coverage**:
- Valid NAT gateway in public subnet
- Invalid NAT gateway with missing name
- Invalid NAT gateway with missing subnet ID
- Invalid NAT gateway with missing allocation ID
- Valid NAT gateway with multiple tags
- NAT gateway-subnet relationship validation
- Elastic IP allocation validation

**Key Constraint**: NAT gateway must be placed in a public subnet.

### 8. `route_table_test.go`
**Purpose**: Tests route table validation and route configurations

**Test Function**: `TestRouteTable`

**Coverage**:
- Valid public route table (0.0.0.0/0 → IGW)
- Valid private route table (0.0.0.0/0 → NAT Gateway)
- Invalid route table with missing name
- Invalid route table with missing VPC ID
- Invalid route with missing destination CIDR
- Invalid route with missing target
- Invalid route with multiple targets
- Route table with local route (edge case)

**Test Scenarios**:
- Public route table: `public-rt` with route `0.0.0.0/0` → IGW
- Private route table: `private-rt` with route `0.0.0.0/0` → NAT Gateway

### 9. `route_table_association_test.go`
**Purpose**: Tests route table associations with subnets

**Test Functions**: 
- `TestRouteTableAssociation` - Standard route table association tests
- `TestSubnetSingleRouteTableAssociation` - Tests subnet single route table constraint

**Coverage**:
- Public route table associated with public subnet
- Private route table associated with three private subnets
- Route table associated with multiple subnets
- Invalid association (route table and subnet in different VPCs)
- **Invalid association (subnet associated with multiple route tables)** - A subnet can only be associated with one route table at a time

**Test Scenarios**:
- **Public Route Table** (`public-rt`):
  - Route: `0.0.0.0/0` → IGW
  - Associated with: Public subnet (10.0.1.0/24)

- **Private Route Table** (`private-rt`):
  - Route: `0.0.0.0/0` → NAT Gateway
  - Associated with: Three private subnets
    - Private subnet 1: `10.0.2.0/24` (us-east-1a)
    - Private subnet 2: `10.0.3.0/24` (us-east-1b)
    - Private subnet 3: `10.0.4.0/24` (us-east-1a)

- **Subnet Single Route Table Constraint** (`TestSubnetSingleRouteTableAssociation`):
  - Tests that a subnet cannot be associated with multiple route tables
  - Scenario: Attempting to associate one subnet with two different route tables
  - Expected: Association with second route table should fail
  - **Key Constraint**: A subnet can only be associated with ONE route table at a time

## Test Output Format

All tests provide detailed console output showing:
- Test case name and description
- Resource details (VPC, subnet, route table, etc.)
- Validation results (✅ PASSED / ❌ FAILED)
- Error messages when validation fails
- Relationship checks (VPC-subnet, route table-subnet, etc.)

### Example Output
```
=== Running test: valid-vpc-in-us-east-1-with-subnets-in-same-region ===
Description: VPC in us-east-1 with subnets in us-east-1a and us-east-1b (same region)
VPC Region: us-east-1

Validating VPC...
✅ PASSED: VPC validation succeeded

Validating 2 subnet(s)...
  Subnet 1: public-subnet (AZ: us-east-1a)
  ✅ PASSED: Subnet validation succeeded
  ✅ PASSED: Subnet AZ us-east-1a is in VPC region us-east-1
  Subnet 2: private-subnet-1 (AZ: us-east-1b)
  ✅ PASSED: Subnet validation succeeded
  ✅ PASSED: Subnet AZ us-east-1b is in VPC region us-east-1
=== Test completed: valid-vpc-in-us-east-1-with-subnets-in-same-region ===
```

## Test Data

The tests use the following consistent test data:

- **Region**: `us-east-1`
- **Availability Zones**: `us-east-1a`, `us-east-1b`
- **VPC CIDR**: `10.0.0.0/16`
- **Public Subnet**: `10.0.1.0/24` (us-east-1a) - Contains NAT gateway
- **Private Subnet 1**: `10.0.2.0/24` (us-east-1a)
- **Private Subnet 2**: `10.0.3.0/24` (us-east-1b)
- **Private Subnet 3**: `10.0.4.0/24` (us-east-1a)

## Dependencies

The tests depend on:
- `github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking` - AWS networking models
- `github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs` - AWS configuration types
- `net` package (standard library) - For CIDR parsing and validation

## Helper Functions

The tests utilize helper functions from the parent package:
- `networking.CIDROverlaps()` - Checks if two CIDR blocks overlap
- `networking.CIDRContains()` - Checks if parent CIDR contains child CIDR
- `Subnet.ValidateCIDRInVPC()` - Validates subnet CIDR is within VPC CIDR
- `Subnet.GetCIDRBlock()` - Parses and returns CIDR block

## Adding New Tests

To add a new test case:

1. Add a new entry to the test cases slice in the appropriate test file
2. Follow the existing table-driven test structure
3. Include a descriptive name and description
4. Set expected error (if any)
5. Add console output for test execution visibility

Example:
```go
{
    name: "new-test-case-name",
    resource: &networking.Resource{
        // ... resource configuration
    },
    expectedError: nil, // or error if expected
    description: "Description of what this test validates",
}
```

## Best Practices

1. **Table-Driven Design**: All tests use table-driven approach for consistency
2. **Descriptive Names**: Test names clearly describe what is being tested
3. **Console Output**: Tests print detailed execution information
4. **Error Validation**: Both positive and negative test cases are included
5. **Relationship Validation**: Tests verify relationships between resources (VPC-subnet, route table-subnet, etc.)

## Troubleshooting

### Tests Failing
- Check that all required fields are set in test data
- Verify CIDR formats are correct (e.g., "10.0.0.0/16")
- Ensure VPC IDs match between related resources
- Check that Availability Zones follow the region-az pattern (e.g., "us-east-1a")

### CIDR Validation Issues
- Ensure subnet CIDR is more specific (larger mask number) than VPC CIDR
- Verify subnet CIDR is within VPC CIDR range
- Check for CIDR overlaps between subnets

### Route Table Issues
- Ensure route has exactly one target (IGW, NAT Gateway, etc.)
- Verify route table and subnets are in the same VPC
- Check that destination CIDR is valid

## Related Documentation

- [AWS Networking Models README](../README.md) - Overview of networking models
- [CIDR Utilities](../cidr_utils.go) - CIDR helper functions
- [Domain Networking Models](../../../../domain/resource/networking/README.md) - Domain layer networking models
