package iam

import (
	"context"
	"errors"
	"fmt"

	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

// IAMService implements AWSIAMService using AWS SDK
type IAMService struct {
	client        *awssdk.AWSClient
	policyService *awssdk.PolicyService
}

// NewIAMService creates a new IAM service implementation
func NewIAMService(client *awssdk.AWSClient) (*IAMService, error) {
	// Initialize policy service for static data fallback
	policyService, err := awssdk.NewPolicyService(client)
	if err != nil {
		// Log warning but don't fail - we can still use SDK directly
		fmt.Printf("Warning: Failed to initialize PolicyService: %v\n", err)
	}

	return &IAMService{
		client:        client,
		policyService: policyService,
	}, nil
}

// Ensure IAMService implements AWSIAMService
var _ AWSIAMService = (*IAMService)(nil)

// Policy operations

func (s *IAMService) CreatePolicy(ctx context.Context, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error) {
	return awssdk.CreatePolicy(ctx, s.client, policy)
}

func (s *IAMService) GetPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error) {
	return awssdk.GetPolicy(ctx, s.client, arn)
}

func (s *IAMService) UpdatePolicy(ctx context.Context, arn string, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error) {
	return awssdk.UpdatePolicy(ctx, s.client, arn, policy)
}

func (s *IAMService) DeletePolicy(ctx context.Context, arn string) error {
	return awssdk.DeletePolicy(ctx, s.client, arn)
}

func (s *IAMService) ListPolicies(ctx context.Context, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error) {
	// List customer managed policies (scope = Local)
	// Don't fetch policy documents for listing (faster performance)
	return awssdk.ListPolicies(ctx, s.client, pathPrefix, types.PolicyScopeTypeLocal, false)
}

func (s *IAMService) ListAWSManagedPolicies(ctx context.Context, scope *string, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error) {
	return awssdk.ListAWSManagedPolicies(ctx, s.client, pathPrefix)
}

func (s *IAMService) GetAWSManagedPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error) {
	return awssdk.GetAWSManagedPolicy(ctx, s.client, arn)
}

// Role operations - Not implemented yet
func (s *IAMService) CreateRole(ctx context.Context, role *awsiam.Role) (*awsoutputs.RoleOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) GetRole(ctx context.Context, name string) (*awsoutputs.RoleOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) UpdateRole(ctx context.Context, name string, role *awsiam.Role) (*awsoutputs.RoleOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DeleteRole(ctx context.Context, name string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListRoles(ctx context.Context, pathPrefix *string) ([]*awsoutputs.RoleOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

// User operations - Not implemented yet
func (s *IAMService) CreateUser(ctx context.Context, user *awsiam.User) (*awsoutputs.UserOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) GetUser(ctx context.Context, name string) (*awsoutputs.UserOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) UpdateUser(ctx context.Context, name string, user *awsiam.User) (*awsoutputs.UserOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DeleteUser(ctx context.Context, name string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListUsers(ctx context.Context, pathPrefix *string) ([]*awsoutputs.UserOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

// Group operations - Not implemented yet
func (s *IAMService) CreateGroup(ctx context.Context, group *awsiam.Group) (*awsoutputs.GroupOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) GetGroup(ctx context.Context, name string) (*awsoutputs.GroupOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) UpdateGroup(ctx context.Context, name string, group *awsiam.Group) (*awsoutputs.GroupOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DeleteGroup(ctx context.Context, name string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListGroups(ctx context.Context, pathPrefix *string) ([]*awsoutputs.GroupOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

// Policy attachment operations - Not implemented yet
func (s *IAMService) AttachPolicyToUser(ctx context.Context, policyARN, userName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DetachPolicyFromUser(ctx context.Context, policyARN, userName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListUserPolicies(ctx context.Context, userName string) ([]*awsoutputs.PolicyOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) AttachPolicyToRole(ctx context.Context, policyARN, roleName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DetachPolicyFromRole(ctx context.Context, policyARN, roleName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListRolePolicies(ctx context.Context, roleName string) ([]*awsoutputs.PolicyOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) AttachPolicyToGroup(ctx context.Context, policyARN, groupName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DetachPolicyFromGroup(ctx context.Context, policyARN, groupName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListGroupPolicies(ctx context.Context, groupName string) ([]*awsoutputs.PolicyOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

// User-Group operations - Not implemented yet
func (s *IAMService) AddUserToGroup(ctx context.Context, userName, groupName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) RemoveUserFromGroup(ctx context.Context, userName, groupName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListGroupUsers(ctx context.Context, groupName string) ([]*awsoutputs.UserOutput, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

// Inline Policy operations - Not implemented yet
func (s *IAMService) PutUserInlinePolicy(ctx context.Context, userName string, policy *awsiam.InlinePolicy) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) GetUserInlinePolicy(ctx context.Context, userName, policyName string) (*awsiam.InlinePolicy, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DeleteUserInlinePolicy(ctx context.Context, userName, policyName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListUserInlinePolicies(ctx context.Context, userName string) ([]string, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) PutRoleInlinePolicy(ctx context.Context, roleName string, policy *awsiam.InlinePolicy) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) GetRoleInlinePolicy(ctx context.Context, roleName, policyName string) (*awsiam.InlinePolicy, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DeleteRoleInlinePolicy(ctx context.Context, roleName, policyName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListRoleInlinePolicies(ctx context.Context, roleName string) ([]string, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) PutGroupInlinePolicy(ctx context.Context, groupName string, policy *awsiam.InlinePolicy) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) GetGroupInlinePolicy(ctx context.Context, groupName, policyName string) (*awsiam.InlinePolicy, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) DeleteGroupInlinePolicy(ctx context.Context, groupName, policyName string) error {
	return errors.New("not implemented: use SDK functions directly")
}

func (s *IAMService) ListGroupInlinePolicies(ctx context.Context, groupName string) ([]string, error) {
	return nil, errors.New("not implemented: use SDK functions directly")
}
