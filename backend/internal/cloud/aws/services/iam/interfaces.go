package iam

import (
	"context"

	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
)

// AWSIAMService defines AWS-specific IAM operations
// This implements cloud provider-specific logic while maintaining domain compatibility
type AWSIAMService interface {
	// Policy operations
	CreatePolicy(ctx context.Context, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error)
	GetPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error)
	UpdatePolicy(ctx context.Context, arn string, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error)
	DeletePolicy(ctx context.Context, arn string) error
	ListPolicies(ctx context.Context, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error)

	// Role operations
	CreateRole(ctx context.Context, role *awsiam.Role) (*awsoutputs.RoleOutput, error)
	GetRole(ctx context.Context, name string) (*awsoutputs.RoleOutput, error)
	UpdateRole(ctx context.Context, name string, role *awsiam.Role) (*awsoutputs.RoleOutput, error)
	DeleteRole(ctx context.Context, name string) error
	ListRoles(ctx context.Context, pathPrefix *string) ([]*awsoutputs.RoleOutput, error)

	// User operations
	CreateUser(ctx context.Context, user *awsiam.User) (*awsoutputs.UserOutput, error)
	GetUser(ctx context.Context, name string) (*awsoutputs.UserOutput, error)
	UpdateUser(ctx context.Context, name string, user *awsiam.User) (*awsoutputs.UserOutput, error)
	DeleteUser(ctx context.Context, name string) error
	ListUsers(ctx context.Context, pathPrefix *string) ([]*awsoutputs.UserOutput, error)

	// Group operations
	CreateGroup(ctx context.Context, group *awsiam.Group) (*awsoutputs.GroupOutput, error)
	GetGroup(ctx context.Context, name string) (*awsoutputs.GroupOutput, error)
	UpdateGroup(ctx context.Context, name string, group *awsiam.Group) (*awsoutputs.GroupOutput, error)
	DeleteGroup(ctx context.Context, name string) error
	ListGroups(ctx context.Context, pathPrefix *string) ([]*awsoutputs.GroupOutput, error)

	// Policy attachment operations
	AttachPolicyToUser(ctx context.Context, policyARN, userName string) error
	DetachPolicyFromUser(ctx context.Context, policyARN, userName string) error
	ListUserPolicies(ctx context.Context, userName string) ([]*awsoutputs.PolicyOutput, error)

	AttachPolicyToRole(ctx context.Context, policyARN, roleName string) error
	DetachPolicyFromRole(ctx context.Context, policyARN, roleName string) error
	ListRolePolicies(ctx context.Context, roleName string) ([]*awsoutputs.PolicyOutput, error)

	AttachPolicyToGroup(ctx context.Context, policyARN, groupName string) error
	DetachPolicyFromGroup(ctx context.Context, policyARN, groupName string) error
	ListGroupPolicies(ctx context.Context, groupName string) ([]*awsoutputs.PolicyOutput, error)

	// User-Group operations
	AddUserToGroup(ctx context.Context, userName, groupName string) error
	RemoveUserFromGroup(ctx context.Context, userName, groupName string) error
	ListGroupUsers(ctx context.Context, groupName string) ([]*awsoutputs.UserOutput, error)

	// Inline Policy operations for Users
	PutUserInlinePolicy(ctx context.Context, userName string, policy *awsiam.InlinePolicy) error
	GetUserInlinePolicy(ctx context.Context, userName, policyName string) (*awsiam.InlinePolicy, error)
	DeleteUserInlinePolicy(ctx context.Context, userName, policyName string) error
	ListUserInlinePolicies(ctx context.Context, userName string) ([]string, error) // Returns policy names

	// Inline Policy operations for Roles
	PutRoleInlinePolicy(ctx context.Context, roleName string, policy *awsiam.InlinePolicy) error
	GetRoleInlinePolicy(ctx context.Context, roleName, policyName string) (*awsiam.InlinePolicy, error)
	DeleteRoleInlinePolicy(ctx context.Context, roleName, policyName string) error
	ListRoleInlinePolicies(ctx context.Context, roleName string) ([]string, error) // Returns policy names

	// Inline Policy operations for Groups
	PutGroupInlinePolicy(ctx context.Context, groupName string, policy *awsiam.InlinePolicy) error
	GetGroupInlinePolicy(ctx context.Context, groupName, policyName string) (*awsiam.InlinePolicy, error)
	DeleteGroupInlinePolicy(ctx context.Context, groupName, policyName string) error
	ListGroupInlinePolicies(ctx context.Context, groupName string) ([]string, error) // Returns policy names

	// AWS Managed Policy operations
	ListAWSManagedPolicies(ctx context.Context, scope *string, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error)
	GetAWSManagedPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error)

	// Instance Profile operations
	CreateInstanceProfile(ctx context.Context, profile *awsiam.InstanceProfile) (*awsoutputs.InstanceProfileOutput, error)
	GetInstanceProfile(ctx context.Context, name string) (*awsoutputs.InstanceProfileOutput, error)
	UpdateInstanceProfile(ctx context.Context, name string, profile *awsiam.InstanceProfile) (*awsoutputs.InstanceProfileOutput, error)
	DeleteInstanceProfile(ctx context.Context, name string) error
	ListInstanceProfiles(ctx context.Context, pathPrefix *string) ([]*awsoutputs.InstanceProfileOutput, error)
	AddRoleToInstanceProfile(ctx context.Context, profileName, roleName string) error
	RemoveRoleFromInstanceProfile(ctx context.Context, profileName, roleName string) error
	GetInstanceProfileRoles(ctx context.Context, profileName string) ([]*awsoutputs.RoleOutput, error)
}
