package iam

// PolicyType represents the type of IAM policy
type PolicyType string

const (
	// PolicyTypeAWSManaged represents AWS managed policies (pre-defined by AWS)
	PolicyTypeAWSManaged PolicyType = "aws_managed"
	// PolicyTypeCustomerManaged represents customer-managed policies (created by users)
	PolicyTypeCustomerManaged PolicyType = "customer_managed"
	// PolicyTypeInline represents inline policies (embedded in users/roles/groups)
	PolicyTypeInline PolicyType = "inline"
)

// PolicyReference represents a reference to a policy with its type
type PolicyReference struct {
	Type     PolicyType // Type of policy
	ARN      *string    // ARN for managed policies (AWS or customer), nil for inline
	Name     string     // Policy name (for inline policies or when ARN is not available)
	Identity string     // Identity name (user/role/group) for inline policies
}
