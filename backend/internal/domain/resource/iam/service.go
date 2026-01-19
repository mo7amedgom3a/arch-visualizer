package iam

import (
	"context"
)

// IAMService defines the interface for IAM resource operations
// This is cloud-agnostic and can be implemented by any cloud provider
type IAMService interface {
	// Policy operations
	CreatePolicy(ctx context.Context, policy *Policy) (*Policy, error)
	GetPolicy(ctx context.Context, arn string) (*Policy, error)
	UpdatePolicy(ctx context.Context, policy *Policy) (*Policy, error)
	DeletePolicy(ctx context.Context, arn string) error
	ListPolicies(ctx context.Context, pathPrefix *string) ([]*Policy, error)

	// Role operations
	CreateRole(ctx context.Context, role *Role) (*Role, error)
	GetRole(ctx context.Context, name string) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) (*Role, error)
	DeleteRole(ctx context.Context, name string) error
	ListRoles(ctx context.Context, pathPrefix *string) ([]*Role, error)

	// User operations
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, name string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, name string) error
	ListUsers(ctx context.Context, pathPrefix *string) ([]*User, error)

	// Group operations
	CreateGroup(ctx context.Context, group *Group) (*Group, error)
	GetGroup(ctx context.Context, name string) (*Group, error)
	UpdateGroup(ctx context.Context, group *Group) (*Group, error)
	DeleteGroup(ctx context.Context, name string) error
	ListGroups(ctx context.Context, pathPrefix *string) ([]*Group, error)

	// Policy attachment operations
	AttachPolicyToUser(ctx context.Context, policyARN, userName string) error
	DetachPolicyFromUser(ctx context.Context, policyARN, userName string) error
	ListUserPolicies(ctx context.Context, userName string) ([]*Policy, error)

	AttachPolicyToRole(ctx context.Context, policyARN, roleName string) error
	DetachPolicyFromRole(ctx context.Context, policyARN, roleName string) error
	ListRolePolicies(ctx context.Context, roleName string) ([]*Policy, error)

	AttachPolicyToGroup(ctx context.Context, policyARN, groupName string) error
	DetachPolicyFromGroup(ctx context.Context, policyARN, groupName string) error
	ListGroupPolicies(ctx context.Context, groupName string) ([]*Policy, error)

	// User-Group operations
	AddUserToGroup(ctx context.Context, userName, groupName string) error
	RemoveUserFromGroup(ctx context.Context, userName, groupName string) error
	ListGroupUsers(ctx context.Context, groupName string) ([]*User, error)

	// Inline Policy operations for Users
	PutUserInlinePolicy(ctx context.Context, userName string, policy *InlinePolicy) error
	GetUserInlinePolicy(ctx context.Context, userName, policyName string) (*InlinePolicy, error)
	DeleteUserInlinePolicy(ctx context.Context, userName, policyName string) error
	ListUserInlinePolicies(ctx context.Context, userName string) ([]*InlinePolicy, error)

	// Inline Policy operations for Roles
	PutRoleInlinePolicy(ctx context.Context, roleName string, policy *InlinePolicy) error
	GetRoleInlinePolicy(ctx context.Context, roleName, policyName string) (*InlinePolicy, error)
	DeleteRoleInlinePolicy(ctx context.Context, roleName, policyName string) error
	ListRoleInlinePolicies(ctx context.Context, roleName string) ([]*InlinePolicy, error)

	// Inline Policy operations for Groups
	PutGroupInlinePolicy(ctx context.Context, groupName string, policy *InlinePolicy) error
	GetGroupInlinePolicy(ctx context.Context, groupName, policyName string) (*InlinePolicy, error)
	DeleteGroupInlinePolicy(ctx context.Context, groupName, policyName string) error
	ListGroupInlinePolicies(ctx context.Context, groupName string) ([]*InlinePolicy, error)

	// AWS Managed Policy operations
	ListAWSManagedPolicies(ctx context.Context, scope *string, pathPrefix *string) ([]*Policy, error)
	GetAWSManagedPolicy(ctx context.Context, arn string) (*Policy, error)
}
