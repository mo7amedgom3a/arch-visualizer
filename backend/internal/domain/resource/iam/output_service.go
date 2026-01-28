package iam

import (
	"context"
)

// IAMOutputService defines the interface for IAM resource operations that return output DTOs
// This is a parallel interface to IAMService, providing output-specific models
type IAMOutputService interface {
	// Policy operations
	CreatePolicyOutput(ctx context.Context, policy *Policy) (*PolicyOutput, error)
	GetPolicyOutput(ctx context.Context, arn string) (*PolicyOutput, error)
	UpdatePolicyOutput(ctx context.Context, policy *Policy) (*PolicyOutput, error)
	ListPoliciesOutput(ctx context.Context, pathPrefix *string) ([]*PolicyOutput, error)

	// Role operations
	CreateRoleOutput(ctx context.Context, role *Role) (*RoleOutput, error)
	GetRoleOutput(ctx context.Context, name string) (*RoleOutput, error)
	UpdateRoleOutput(ctx context.Context, role *Role) (*RoleOutput, error)
	ListRolesOutput(ctx context.Context, pathPrefix *string) ([]*RoleOutput, error)

	// User operations
	CreateUserOutput(ctx context.Context, user *User) (*UserOutput, error)
	GetUserOutput(ctx context.Context, name string) (*UserOutput, error)
	UpdateUserOutput(ctx context.Context, user *User) (*UserOutput, error)
	ListUsersOutput(ctx context.Context, pathPrefix *string) ([]*UserOutput, error)

	// Group operations
	CreateGroupOutput(ctx context.Context, group *Group) (*GroupOutput, error)
	GetGroupOutput(ctx context.Context, name string) (*GroupOutput, error)
	UpdateGroupOutput(ctx context.Context, group *Group) (*GroupOutput, error)
	ListGroupsOutput(ctx context.Context, pathPrefix *string) ([]*GroupOutput, error)

	// Policy attachment operations (these don't return outputs, but kept for interface consistency)
	AttachPolicyToUser(ctx context.Context, policyARN, userName string) error
	DetachPolicyFromUser(ctx context.Context, policyARN, userName string) error
	ListUserPoliciesOutput(ctx context.Context, userName string) ([]*PolicyOutput, error)

	AttachPolicyToRole(ctx context.Context, policyARN, roleName string) error
	DetachPolicyFromRole(ctx context.Context, policyARN, roleName string) error
	ListRolePoliciesOutput(ctx context.Context, roleName string) ([]*PolicyOutput, error)

	AttachPolicyToGroup(ctx context.Context, policyARN, groupName string) error
	DetachPolicyFromGroup(ctx context.Context, policyARN, groupName string) error
	ListGroupPoliciesOutput(ctx context.Context, groupName string) ([]*PolicyOutput, error)

	// User-Group operations
	AddUserToGroup(ctx context.Context, userName, groupName string) error
	RemoveUserFromGroup(ctx context.Context, userName, groupName string) error
	ListGroupUsersOutput(ctx context.Context, groupName string) ([]*UserOutput, error)

	// Inline Policy operations for Users
	PutUserInlinePolicy(ctx context.Context, userName string, policy *InlinePolicy) error
	GetUserInlinePolicyOutput(ctx context.Context, userName, policyName string) (*InlinePolicyOutput, error)
	DeleteUserInlinePolicy(ctx context.Context, userName, policyName string) error
	ListUserInlinePoliciesOutput(ctx context.Context, userName string) ([]*InlinePolicyOutput, error)

	// Inline Policy operations for Roles
	PutRoleInlinePolicy(ctx context.Context, roleName string, policy *InlinePolicy) error
	GetRoleInlinePolicyOutput(ctx context.Context, roleName, policyName string) (*InlinePolicyOutput, error)
	DeleteRoleInlinePolicy(ctx context.Context, roleName, policyName string) error
	ListRoleInlinePoliciesOutput(ctx context.Context, roleName string) ([]*InlinePolicyOutput, error)

	// Inline Policy operations for Groups
	PutGroupInlinePolicy(ctx context.Context, groupName string, policy *InlinePolicy) error
	GetGroupInlinePolicyOutput(ctx context.Context, groupName, policyName string) (*InlinePolicyOutput, error)
	DeleteGroupInlinePolicy(ctx context.Context, groupName, policyName string) error
	ListGroupInlinePoliciesOutput(ctx context.Context, groupName string) ([]*InlinePolicyOutput, error)

	// AWS Managed Policy operations
	ListAWSManagedPoliciesOutput(ctx context.Context, scope *string, pathPrefix *string) ([]*PolicyOutput, error)
	GetAWSManagedPolicyOutput(ctx context.Context, arn string) (*PolicyOutput, error)

	// Instance Profile operations
	CreateInstanceProfileOutput(ctx context.Context, profile *InstanceProfile) (*InstanceProfileOutput, error)
	GetInstanceProfileOutput(ctx context.Context, name string) (*InstanceProfileOutput, error)
	UpdateInstanceProfileOutput(ctx context.Context, profile *InstanceProfile) (*InstanceProfileOutput, error)
	ListInstanceProfilesOutput(ctx context.Context, pathPrefix *string) ([]*InstanceProfileOutput, error)
	AddRoleToInstanceProfile(ctx context.Context, profileName, roleName string) error
	RemoveRoleFromInstanceProfile(ctx context.Context, profileName, roleName string) error
	GetInstanceProfileRolesOutput(ctx context.Context, profileName string) ([]*RoleOutput, error)
}
