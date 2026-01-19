package iam

import (
	"context"
	"fmt"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/iam"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/iam"
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
)

// AWSIAMAdapter adapts AWS-specific IAM service to domain IAM service
// This implements the Adapter pattern, allowing the domain layer to work with cloud-specific implementations
type AWSIAMAdapter struct {
	awsService awsservice.AWSIAMService
}

// NewAWSIAMAdapter creates a new AWS IAM adapter
func NewAWSIAMAdapter(awsService awsservice.AWSIAMService) domainiam.IAMService {
	return &AWSIAMAdapter{
		awsService: awsService,
	}
}

// Ensure AWSIAMAdapter implements IAMService
var _ domainiam.IAMService = (*AWSIAMAdapter)(nil)

// Policy Operations

func (a *AWSIAMAdapter) CreatePolicy(ctx context.Context, policy *domainiam.Policy) (*domainiam.Policy, error) {
	if err := policy.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsPolicy := awsmapper.FromDomainPolicy(policy)
	if err := awsPolicy.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsPolicyOutput, err := a.awsService.CreatePolicy(ctx, awsPolicy)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput), nil
}

func (a *AWSIAMAdapter) GetPolicy(ctx context.Context, arn string) (*domainiam.Policy, error) {
	awsPolicyOutput, err := a.awsService.GetPolicy(ctx, arn)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput), nil
}

func (a *AWSIAMAdapter) UpdatePolicy(ctx context.Context, policy *domainiam.Policy) (*domainiam.Policy, error) {
	if err := policy.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	if policy.ARN == nil || *policy.ARN == "" {
		return nil, fmt.Errorf("policy ARN is required for update")
	}

	awsPolicy := awsmapper.FromDomainPolicy(policy)
	if err := awsPolicy.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsPolicyOutput, err := a.awsService.UpdatePolicy(ctx, *policy.ARN, awsPolicy)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput), nil
}

func (a *AWSIAMAdapter) DeletePolicy(ctx context.Context, arn string) error {
	if err := a.awsService.DeletePolicy(ctx, arn); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListPolicies(ctx context.Context, pathPrefix *string) ([]*domainiam.Policy, error) {
	awsPolicyOutputs, err := a.awsService.ListPolicies(ctx, pathPrefix)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainPolicies := make([]*domainiam.Policy, len(awsPolicyOutputs))
	for i, awsPolicyOutput := range awsPolicyOutputs {
		domainPolicies[i] = awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput)
	}

	return domainPolicies, nil
}

// Role Operations

func (a *AWSIAMAdapter) CreateRole(ctx context.Context, role *domainiam.Role) (*domainiam.Role, error) {
	if err := role.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsRole := awsmapper.FromDomainRole(role)
	if err := awsRole.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsRoleOutput, err := a.awsService.CreateRole(ctx, awsRole)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainRoleFromOutput(awsRoleOutput), nil
}

func (a *AWSIAMAdapter) GetRole(ctx context.Context, name string) (*domainiam.Role, error) {
	awsRoleOutput, err := a.awsService.GetRole(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainRoleFromOutput(awsRoleOutput), nil
}

func (a *AWSIAMAdapter) UpdateRole(ctx context.Context, role *domainiam.Role) (*domainiam.Role, error) {
	if err := role.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	if role.Name == "" {
		return nil, fmt.Errorf("role name is required for update")
	}

	awsRole := awsmapper.FromDomainRole(role)
	if err := awsRole.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsRoleOutput, err := a.awsService.UpdateRole(ctx, role.Name, awsRole)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainRoleFromOutput(awsRoleOutput), nil
}

func (a *AWSIAMAdapter) DeleteRole(ctx context.Context, name string) error {
	if err := a.awsService.DeleteRole(ctx, name); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListRoles(ctx context.Context, pathPrefix *string) ([]*domainiam.Role, error) {
	awsRoleOutputs, err := a.awsService.ListRoles(ctx, pathPrefix)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainRoles := make([]*domainiam.Role, len(awsRoleOutputs))
	for i, awsRoleOutput := range awsRoleOutputs {
		domainRoles[i] = awsmapper.ToDomainRoleFromOutput(awsRoleOutput)
	}

	return domainRoles, nil
}

// User Operations

func (a *AWSIAMAdapter) CreateUser(ctx context.Context, user *domainiam.User) (*domainiam.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsUser := awsmapper.FromDomainUser(user)
	if err := awsUser.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsUserOutput, err := a.awsService.CreateUser(ctx, awsUser)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainUserFromOutput(awsUserOutput), nil
}

func (a *AWSIAMAdapter) GetUser(ctx context.Context, name string) (*domainiam.User, error) {
	awsUserOutput, err := a.awsService.GetUser(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainUserFromOutput(awsUserOutput), nil
}

func (a *AWSIAMAdapter) UpdateUser(ctx context.Context, user *domainiam.User) (*domainiam.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	if user.Name == "" {
		return nil, fmt.Errorf("user name is required for update")
	}

	awsUser := awsmapper.FromDomainUser(user)
	if err := awsUser.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsUserOutput, err := a.awsService.UpdateUser(ctx, user.Name, awsUser)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainUserFromOutput(awsUserOutput), nil
}

func (a *AWSIAMAdapter) DeleteUser(ctx context.Context, name string) error {
	if err := a.awsService.DeleteUser(ctx, name); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListUsers(ctx context.Context, pathPrefix *string) ([]*domainiam.User, error) {
	awsUserOutputs, err := a.awsService.ListUsers(ctx, pathPrefix)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainUsers := make([]*domainiam.User, len(awsUserOutputs))
	for i, awsUserOutput := range awsUserOutputs {
		domainUsers[i] = awsmapper.ToDomainUserFromOutput(awsUserOutput)
	}

	return domainUsers, nil
}

// Group Operations

func (a *AWSIAMAdapter) CreateGroup(ctx context.Context, group *domainiam.Group) (*domainiam.Group, error) {
	if err := group.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsGroup := awsmapper.FromDomainGroup(group)
	if err := awsGroup.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsGroupOutput, err := a.awsService.CreateGroup(ctx, awsGroup)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainGroupFromOutput(awsGroupOutput), nil
}

func (a *AWSIAMAdapter) GetGroup(ctx context.Context, name string) (*domainiam.Group, error) {
	awsGroupOutput, err := a.awsService.GetGroup(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainGroupFromOutput(awsGroupOutput), nil
}

func (a *AWSIAMAdapter) UpdateGroup(ctx context.Context, group *domainiam.Group) (*domainiam.Group, error) {
	if err := group.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	if group.Name == "" {
		return nil, fmt.Errorf("group name is required for update")
	}

	awsGroup := awsmapper.FromDomainGroup(group)
	if err := awsGroup.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsGroupOutput, err := a.awsService.UpdateGroup(ctx, group.Name, awsGroup)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainGroupFromOutput(awsGroupOutput), nil
}

func (a *AWSIAMAdapter) DeleteGroup(ctx context.Context, name string) error {
	if err := a.awsService.DeleteGroup(ctx, name); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListGroups(ctx context.Context, pathPrefix *string) ([]*domainiam.Group, error) {
	awsGroupOutputs, err := a.awsService.ListGroups(ctx, pathPrefix)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainGroups := make([]*domainiam.Group, len(awsGroupOutputs))
	for i, awsGroupOutput := range awsGroupOutputs {
		domainGroups[i] = awsmapper.ToDomainGroupFromOutput(awsGroupOutput)
	}

	return domainGroups, nil
}

// Policy Attachment Operations

func (a *AWSIAMAdapter) AttachPolicyToUser(ctx context.Context, policyARN, userName string) error {
	if err := a.awsService.AttachPolicyToUser(ctx, policyARN, userName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) DetachPolicyFromUser(ctx context.Context, policyARN, userName string) error {
	if err := a.awsService.DetachPolicyFromUser(ctx, policyARN, userName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListUserPolicies(ctx context.Context, userName string) ([]*domainiam.Policy, error) {
	awsPolicyOutputs, err := a.awsService.ListUserPolicies(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainPolicies := make([]*domainiam.Policy, len(awsPolicyOutputs))
	for i, awsPolicyOutput := range awsPolicyOutputs {
		domainPolicies[i] = awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput)
	}

	return domainPolicies, nil
}

func (a *AWSIAMAdapter) AttachPolicyToRole(ctx context.Context, policyARN, roleName string) error {
	if err := a.awsService.AttachPolicyToRole(ctx, policyARN, roleName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) DetachPolicyFromRole(ctx context.Context, policyARN, roleName string) error {
	if err := a.awsService.DetachPolicyFromRole(ctx, policyARN, roleName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListRolePolicies(ctx context.Context, roleName string) ([]*domainiam.Policy, error) {
	awsPolicyOutputs, err := a.awsService.ListRolePolicies(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainPolicies := make([]*domainiam.Policy, len(awsPolicyOutputs))
	for i, awsPolicyOutput := range awsPolicyOutputs {
		domainPolicies[i] = awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput)
	}

	return domainPolicies, nil
}

func (a *AWSIAMAdapter) AttachPolicyToGroup(ctx context.Context, policyARN, groupName string) error {
	if err := a.awsService.AttachPolicyToGroup(ctx, policyARN, groupName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) DetachPolicyFromGroup(ctx context.Context, policyARN, groupName string) error {
	if err := a.awsService.DetachPolicyFromGroup(ctx, policyARN, groupName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListGroupPolicies(ctx context.Context, groupName string) ([]*domainiam.Policy, error) {
	awsPolicyOutputs, err := a.awsService.ListGroupPolicies(ctx, groupName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainPolicies := make([]*domainiam.Policy, len(awsPolicyOutputs))
	for i, awsPolicyOutput := range awsPolicyOutputs {
		domainPolicies[i] = awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput)
	}

	return domainPolicies, nil
}

// User-Group Operations

func (a *AWSIAMAdapter) AddUserToGroup(ctx context.Context, userName, groupName string) error {
	if err := a.awsService.AddUserToGroup(ctx, userName, groupName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) RemoveUserFromGroup(ctx context.Context, userName, groupName string) error {
	if err := a.awsService.RemoveUserFromGroup(ctx, userName, groupName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListGroupUsers(ctx context.Context, groupName string) ([]*domainiam.User, error) {
	awsUserOutputs, err := a.awsService.ListGroupUsers(ctx, groupName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainUsers := make([]*domainiam.User, len(awsUserOutputs))
	for i, awsUserOutput := range awsUserOutputs {
		domainUsers[i] = awsmapper.ToDomainUserFromOutput(awsUserOutput)
	}

	return domainUsers, nil
}

// Inline Policy Operations for Users

func (a *AWSIAMAdapter) PutUserInlinePolicy(ctx context.Context, userName string, policy *domainiam.InlinePolicy) error {
	if err := policy.Validate(); err != nil {
		return fmt.Errorf("domain validation failed: %w", err)
	}

	awsInlinePolicy := awsmapper.FromDomainInlinePolicy(policy)
	if err := a.awsService.PutUserInlinePolicy(ctx, userName, awsInlinePolicy); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) GetUserInlinePolicy(ctx context.Context, userName, policyName string) (*domainiam.InlinePolicy, error) {
	awsInlinePolicy, err := a.awsService.GetUserInlinePolicy(ctx, userName, policyName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return awsmapper.ToDomainInlinePolicy(awsInlinePolicy), nil
}

func (a *AWSIAMAdapter) DeleteUserInlinePolicy(ctx context.Context, userName, policyName string) error {
	if err := a.awsService.DeleteUserInlinePolicy(ctx, userName, policyName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListUserInlinePolicies(ctx context.Context, userName string) ([]*domainiam.InlinePolicy, error) {
	policyNames, err := a.awsService.ListUserInlinePolicies(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	policies := make([]*domainiam.InlinePolicy, 0, len(policyNames))
	for _, policyName := range policyNames {
		policy, err := a.GetUserInlinePolicy(ctx, userName, policyName)
		if err != nil {
			return nil, fmt.Errorf("failed to get inline policy %s: %w", policyName, err)
		}
		policies = append(policies, policy)
	}
	return policies, nil
}

// Inline Policy Operations for Roles

func (a *AWSIAMAdapter) PutRoleInlinePolicy(ctx context.Context, roleName string, policy *domainiam.InlinePolicy) error {
	if err := policy.Validate(); err != nil {
		return fmt.Errorf("domain validation failed: %w", err)
	}

	awsInlinePolicy := awsmapper.FromDomainInlinePolicy(policy)
	if err := a.awsService.PutRoleInlinePolicy(ctx, roleName, awsInlinePolicy); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) GetRoleInlinePolicy(ctx context.Context, roleName, policyName string) (*domainiam.InlinePolicy, error) {
	awsInlinePolicy, err := a.awsService.GetRoleInlinePolicy(ctx, roleName, policyName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return awsmapper.ToDomainInlinePolicy(awsInlinePolicy), nil
}

func (a *AWSIAMAdapter) DeleteRoleInlinePolicy(ctx context.Context, roleName, policyName string) error {
	if err := a.awsService.DeleteRoleInlinePolicy(ctx, roleName, policyName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListRoleInlinePolicies(ctx context.Context, roleName string) ([]*domainiam.InlinePolicy, error) {
	policyNames, err := a.awsService.ListRoleInlinePolicies(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	policies := make([]*domainiam.InlinePolicy, 0, len(policyNames))
	for _, policyName := range policyNames {
		policy, err := a.GetRoleInlinePolicy(ctx, roleName, policyName)
		if err != nil {
			return nil, fmt.Errorf("failed to get inline policy %s: %w", policyName, err)
		}
		policies = append(policies, policy)
	}
	return policies, nil
}

// Inline Policy Operations for Groups

func (a *AWSIAMAdapter) PutGroupInlinePolicy(ctx context.Context, groupName string, policy *domainiam.InlinePolicy) error {
	if err := policy.Validate(); err != nil {
		return fmt.Errorf("domain validation failed: %w", err)
	}

	awsInlinePolicy := awsmapper.FromDomainInlinePolicy(policy)
	if err := a.awsService.PutGroupInlinePolicy(ctx, groupName, awsInlinePolicy); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) GetGroupInlinePolicy(ctx context.Context, groupName, policyName string) (*domainiam.InlinePolicy, error) {
	awsInlinePolicy, err := a.awsService.GetGroupInlinePolicy(ctx, groupName, policyName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return awsmapper.ToDomainInlinePolicy(awsInlinePolicy), nil
}

func (a *AWSIAMAdapter) DeleteGroupInlinePolicy(ctx context.Context, groupName, policyName string) error {
	if err := a.awsService.DeleteGroupInlinePolicy(ctx, groupName, policyName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListGroupInlinePolicies(ctx context.Context, groupName string) ([]*domainiam.InlinePolicy, error) {
	policyNames, err := a.awsService.ListGroupInlinePolicies(ctx, groupName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	policies := make([]*domainiam.InlinePolicy, 0, len(policyNames))
	for _, policyName := range policyNames {
		policy, err := a.GetGroupInlinePolicy(ctx, groupName, policyName)
		if err != nil {
			return nil, fmt.Errorf("failed to get inline policy %s: %w", policyName, err)
		}
		policies = append(policies, policy)
	}
	return policies, nil
}

// AWS Managed Policy Operations

func (a *AWSIAMAdapter) ListAWSManagedPolicies(ctx context.Context, scope *string, pathPrefix *string) ([]*domainiam.Policy, error) {
	awsPolicyOutputs, err := a.awsService.ListAWSManagedPolicies(ctx, scope, pathPrefix)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainPolicies := make([]*domainiam.Policy, len(awsPolicyOutputs))
	for i, awsPolicyOutput := range awsPolicyOutputs {
		domainPolicies[i] = awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput)
	}

	return domainPolicies, nil
}

func (a *AWSIAMAdapter) GetAWSManagedPolicy(ctx context.Context, arn string) (*domainiam.Policy, error) {
	awsPolicyOutput, err := a.awsService.GetAWSManagedPolicy(ctx, arn)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return awsmapper.ToDomainPolicyFromOutput(awsPolicyOutput), nil
}

// Instance Profile operations

func (a *AWSIAMAdapter) CreateInstanceProfile(ctx context.Context, profile *domainiam.InstanceProfile) (*domainiam.InstanceProfile, error) {
	if err := profile.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsProfile := awsmapper.FromDomainInstanceProfile(profile)
	if err := awsProfile.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsProfileOutput, err := a.awsService.CreateInstanceProfile(ctx, awsProfile)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceProfileFromOutput(awsProfileOutput), nil
}

func (a *AWSIAMAdapter) GetInstanceProfile(ctx context.Context, name string) (*domainiam.InstanceProfile, error) {
	awsProfileOutput, err := a.awsService.GetInstanceProfile(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceProfileFromOutput(awsProfileOutput), nil
}

func (a *AWSIAMAdapter) UpdateInstanceProfile(ctx context.Context, profile *domainiam.InstanceProfile) (*domainiam.InstanceProfile, error) {
	if err := profile.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}
	if profile.Name == "" {
		return nil, fmt.Errorf("instance profile name is required for update")
	}

	awsProfile := awsmapper.FromDomainInstanceProfile(profile)
	if err := awsProfile.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsProfileOutput, err := a.awsService.UpdateInstanceProfile(ctx, profile.Name, awsProfile)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceProfileFromOutput(awsProfileOutput), nil
}

func (a *AWSIAMAdapter) DeleteInstanceProfile(ctx context.Context, name string) error {
	if err := a.awsService.DeleteInstanceProfile(ctx, name); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) ListInstanceProfiles(ctx context.Context, pathPrefix *string) ([]*domainiam.InstanceProfile, error) {
	awsProfileOutputs, err := a.awsService.ListInstanceProfiles(ctx, pathPrefix)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainProfiles := make([]*domainiam.InstanceProfile, len(awsProfileOutputs))
	for i, awsProfileOutput := range awsProfileOutputs {
		domainProfiles[i] = awsmapper.ToDomainInstanceProfileFromOutput(awsProfileOutput)
	}
	return domainProfiles, nil
}

func (a *AWSIAMAdapter) AddRoleToInstanceProfile(ctx context.Context, profileName, roleName string) error {
	if err := a.awsService.AddRoleToInstanceProfile(ctx, profileName, roleName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) RemoveRoleFromInstanceProfile(ctx context.Context, profileName, roleName string) error {
	if err := a.awsService.RemoveRoleFromInstanceProfile(ctx, profileName, roleName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSIAMAdapter) GetInstanceProfileRoles(ctx context.Context, profileName string) ([]*domainiam.Role, error) {
	awsRoleOutputs, err := a.awsService.GetInstanceProfileRoles(ctx, profileName)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainRoles := make([]*domainiam.Role, len(awsRoleOutputs))
	for i, awsRoleOutput := range awsRoleOutputs {
		domainRoles[i] = awsmapper.ToDomainRoleFromOutput(awsRoleOutput)
	}
	return domainRoles, nil
}
