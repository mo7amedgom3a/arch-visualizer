package iam

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
)

// IAMService implements AWSIAMService with deterministic virtual operations
type IAMService struct {
	policyRepo *PolicyRepository
}

// NewIAMService creates a new IAM service implementation
func NewIAMService() *IAMService {
	repo := NewPolicyRepository()
	// Best effort load, if it fails logs will show in simulation/server logs
	// In a real app we might want to panic or return error, but here we keep it simple
	_ = repo.LoadPolicies()
	return &IAMService{
		policyRepo: repo,
	}
}

// Ensure IAMService implements AWSIAMService
var _ AWSIAMService = (*IAMService)(nil)

// Policy operations

func (s *IAMService) CreatePolicy(ctx context.Context, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error) {
	if policy == nil {
		return nil, fmt.Errorf("policy is nil")
	}

	path := "/"
	if policy.Path != nil {
		path = *policy.Path
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:policy%s%s", path, policy.Name)
	if path == "" {
		arn = fmt.Sprintf("arn:aws:iam::123456789012:policy/%s", policy.Name)
	}

	return &awsoutputs.PolicyOutput{
		ARN:              arn,
		ID:               arn,
		Name:             policy.Name,
		Description:      policy.Description,
		Path:             path,
		PolicyDocument:   policy.PolicyDocument,
		CreateDate:       services.GetFixedTimestamp(),
		UpdateDate:       services.GetFixedTimestamp(),
		DefaultVersionID: services.StringPtr("v1"),
		AttachmentCount:  0,
		IsAttachable:     true,
		Tags:             policy.Tags,
		IsAWSManaged:     false,
	}, nil
}

func (s *IAMService) GetPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error) {
	return &awsoutputs.PolicyOutput{
		ARN:              arn,
		ID:               arn,
		Name:             "test-policy",
		Description:      services.StringPtr("Test policy"),
		Path:             "/",
		PolicyDocument:   `{"Version":"2012-10-17","Statement":[]}`,
		CreateDate:       services.GetFixedTimestamp(),
		UpdateDate:       services.GetFixedTimestamp(),
		DefaultVersionID: services.StringPtr("v1"),
		AttachmentCount:  0,
		IsAttachable:     true,
		Tags:             []configs.Tag{},
		IsAWSManaged:     false,
	}, nil
}

func (s *IAMService) UpdatePolicy(ctx context.Context, arn string, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error) {
	return s.CreatePolicy(ctx, policy)
}

func (s *IAMService) DeletePolicy(ctx context.Context, arn string) error {
	return nil
}

func (s *IAMService) ListPolicies(ctx context.Context, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error) {
	path := "/"
	if pathPrefix != nil {
		path = *pathPrefix
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:policy%stest-policy", path)
	if path == "" {
		arn = "arn:aws:iam::123456789012:policy/test-policy"
	}

	return []*awsoutputs.PolicyOutput{
		{
			ARN:              arn,
			ID:               arn,
			Name:             "test-policy",
			Description:      services.StringPtr("Test policy"),
			Path:             path,
			PolicyDocument:   `{"Version":"2012-10-17","Statement":[]}`,
			CreateDate:       services.GetFixedTimestamp(),
			UpdateDate:       services.GetFixedTimestamp(),
			DefaultVersionID: services.StringPtr("v1"),
			AttachmentCount:  0,
			IsAttachable:     true,
			Tags:             []configs.Tag{},
			IsAWSManaged:     false,
		},
	}, nil
}

func (s *IAMService) ListAWSManagedPolicies(ctx context.Context, scope *string, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error) {
	// Use repository to list policies
	// If scope is provided (e.g. "S3", "Lambda"), use it as filter
	filter := ""
	if scope != nil {
		filter = *scope
	}
	return s.policyRepo.ListPolicies(filter), nil
}

func (s *IAMService) GetAWSManagedPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error) {
	return &awsoutputs.PolicyOutput{
		ARN:              arn,
		ID:               arn,
		Name:             "ReadOnlyAccess",
		Description:      services.StringPtr("Provides read-only access to AWS services and resources"),
		Path:             "/",
		PolicyDocument:   `{"Version":"2012-10-17","Statement":[]}`,
		CreateDate:       services.GetFixedTimestamp(),
		UpdateDate:       services.GetFixedTimestamp(),
		DefaultVersionID: services.StringPtr("v1"),
		AttachmentCount:  0,
		IsAttachable:     true,
		Tags:             []configs.Tag{},
		IsAWSManaged:     true,
	}, nil
}

// Role operations

func (s *IAMService) CreateRole(ctx context.Context, role *awsiam.Role) (*awsoutputs.RoleOutput, error) {
	if role == nil {
		return nil, fmt.Errorf("role is nil")
	}

	path := "/"
	if role.Path != nil {
		path = *role.Path
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:role%s%s", path, role.Name)
	if path == "" {
		arn = fmt.Sprintf("arn:aws:iam::123456789012:role/%s", role.Name)
	}
	uniqueID := services.GenerateDeterministicID(role.Name)

	return &awsoutputs.RoleOutput{
		ARN:                 arn,
		ID:                  role.Name,
		Name:                role.Name,
		UniqueID:            uniqueID,
		Description:         role.Description,
		Path:                path,
		AssumeRolePolicy:    role.AssumeRolePolicy,
		PermissionsBoundary: role.PermissionsBoundary,
		CreateDate:          services.GetFixedTimestamp(),
		MaxSessionDuration:  nil,
		Tags:                role.Tags,
		IsVirtual:           role.IsVirtual,
	}, nil
}

func (s *IAMService) GetRole(ctx context.Context, name string) (*awsoutputs.RoleOutput, error) {
	arn := fmt.Sprintf("arn:aws:iam::123456789012:role/%s", name)
	uniqueID := services.GenerateDeterministicID(name)

	return &awsoutputs.RoleOutput{
		ARN:                 arn,
		ID:                  name,
		Name:                name,
		UniqueID:            uniqueID,
		Description:         services.StringPtr("Test role"),
		Path:                "/",
		AssumeRolePolicy:    `{"Version":"2012-10-17","Statement":[]}`,
		PermissionsBoundary: nil,
		CreateDate:          services.GetFixedTimestamp(),
		MaxSessionDuration:  nil,
		Tags:                []configs.Tag{},
	}, nil
}

func (s *IAMService) UpdateRole(ctx context.Context, name string, role *awsiam.Role) (*awsoutputs.RoleOutput, error) {
	return s.CreateRole(ctx, role)
}

func (s *IAMService) DeleteRole(ctx context.Context, name string) error {
	return nil
}

func (s *IAMService) ListRoles(ctx context.Context, pathPrefix *string) ([]*awsoutputs.RoleOutput, error) {
	path := "/"
	if pathPrefix != nil {
		path = *pathPrefix
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:role%stest-role", path)
	if path == "" {
		arn = "arn:aws:iam::123456789012:role/test-role"
	}
	uniqueID := services.GenerateDeterministicID("test-role")

	return []*awsoutputs.RoleOutput{
		{
			ARN:                 arn,
			ID:                  "test-role",
			Name:                "test-role",
			UniqueID:            uniqueID,
			Description:         services.StringPtr("Test role"),
			Path:                path,
			AssumeRolePolicy:    `{"Version":"2012-10-17","Statement":[]}`,
			PermissionsBoundary: nil,
			CreateDate:          services.GetFixedTimestamp(),
			MaxSessionDuration:  nil,
			Tags:                []configs.Tag{},
		},
	}, nil
}

// User operations

func (s *IAMService) CreateUser(ctx context.Context, user *awsiam.User) (*awsoutputs.UserOutput, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	path := "/"
	if user.Path != nil {
		path = *user.Path
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:user%s%s", path, user.Name)
	if path == "" {
		arn = fmt.Sprintf("arn:aws:iam::123456789012:user/%s", user.Name)
	}
	uniqueID := services.GenerateDeterministicID(user.Name)

	return &awsoutputs.UserOutput{
		ARN:                 arn,
		ID:                  user.Name,
		Name:                user.Name,
		UniqueID:            uniqueID,
		Path:                path,
		PermissionsBoundary: user.PermissionsBoundary,
		CreateDate:          services.GetFixedTimestamp(),
		PasswordLastUsed:    nil,
		Tags:                user.Tags,
		IsVirtual:           user.IsVirtual,
	}, nil
}

func (s *IAMService) GetUser(ctx context.Context, name string) (*awsoutputs.UserOutput, error) {
	arn := fmt.Sprintf("arn:aws:iam::123456789012:user/%s", name)
	uniqueID := services.GenerateDeterministicID(name)

	return &awsoutputs.UserOutput{
		ARN:                 arn,
		ID:                  name,
		Name:                name,
		UniqueID:            uniqueID,
		Path:                "/",
		PermissionsBoundary: nil,
		CreateDate:          services.GetFixedTimestamp(),
		PasswordLastUsed:    nil,
		Tags:                []configs.Tag{},
	}, nil
}

func (s *IAMService) UpdateUser(ctx context.Context, name string, user *awsiam.User) (*awsoutputs.UserOutput, error) {
	return s.CreateUser(ctx, user)
}

func (s *IAMService) DeleteUser(ctx context.Context, name string) error {
	return nil
}

func (s *IAMService) ListUsers(ctx context.Context, pathPrefix *string) ([]*awsoutputs.UserOutput, error) {
	path := "/"
	if pathPrefix != nil {
		path = *pathPrefix
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:user%stest-user", path)
	if path == "" {
		arn = "arn:aws:iam::123456789012:user/test-user"
	}
	uniqueID := services.GenerateDeterministicID("test-user")

	return []*awsoutputs.UserOutput{
		{
			ARN:                 arn,
			ID:                  "test-user",
			Name:                "test-user",
			UniqueID:            uniqueID,
			Path:                path,
			PermissionsBoundary: nil,
			CreateDate:          services.GetFixedTimestamp(),
			PasswordLastUsed:    nil,
			Tags:                []configs.Tag{},
		},
	}, nil
}

// Group operations

func (s *IAMService) CreateGroup(ctx context.Context, group *awsiam.Group) (*awsoutputs.GroupOutput, error) {
	if group == nil {
		return nil, fmt.Errorf("group is nil")
	}

	path := "/"
	if group.Path != nil {
		path = *group.Path
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:group%s%s", path, group.Name)
	if path == "" {
		arn = fmt.Sprintf("arn:aws:iam::123456789012:group/%s", group.Name)
	}
	uniqueID := services.GenerateDeterministicID(group.Name)

	return &awsoutputs.GroupOutput{
		ARN:        arn,
		ID:         group.Name,
		Name:       group.Name,
		UniqueID:   uniqueID,
		Path:       path,
		CreateDate: services.GetFixedTimestamp(),
		Tags:       group.Tags,
	}, nil
}

func (s *IAMService) GetGroup(ctx context.Context, name string) (*awsoutputs.GroupOutput, error) {
	arn := fmt.Sprintf("arn:aws:iam::123456789012:group/%s", name)
	uniqueID := services.GenerateDeterministicID(name)

	return &awsoutputs.GroupOutput{
		ARN:        arn,
		ID:         name,
		Name:       name,
		UniqueID:   uniqueID,
		Path:       "/",
		CreateDate: services.GetFixedTimestamp(),
		Tags:       []configs.Tag{},
	}, nil
}

func (s *IAMService) UpdateGroup(ctx context.Context, name string, group *awsiam.Group) (*awsoutputs.GroupOutput, error) {
	return s.CreateGroup(ctx, group)
}

func (s *IAMService) DeleteGroup(ctx context.Context, name string) error {
	return nil
}

func (s *IAMService) ListGroups(ctx context.Context, pathPrefix *string) ([]*awsoutputs.GroupOutput, error) {
	path := "/"
	if pathPrefix != nil {
		path = *pathPrefix
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:group%stest-group", path)
	if path == "" {
		arn = "arn:aws:iam::123456789012:group/test-group"
	}
	uniqueID := services.GenerateDeterministicID("test-group")

	return []*awsoutputs.GroupOutput{
		{
			ARN:        arn,
			ID:         "test-group",
			Name:       "test-group",
			UniqueID:   uniqueID,
			Path:       path,
			CreateDate: services.GetFixedTimestamp(),
			Tags:       []configs.Tag{},
		},
	}, nil
}

// Policy attachment operations

func (s *IAMService) AttachPolicyToUser(ctx context.Context, policyARN, userName string) error {
	return nil
}

func (s *IAMService) DetachPolicyFromUser(ctx context.Context, policyARN, userName string) error {
	return nil
}

func (s *IAMService) ListUserPolicies(ctx context.Context, userName string) ([]*awsoutputs.PolicyOutput, error) {
	return []*awsoutputs.PolicyOutput{}, nil
}

func (s *IAMService) AttachPolicyToRole(ctx context.Context, policyARN, roleName string) error {
	return nil
}

func (s *IAMService) DetachPolicyFromRole(ctx context.Context, policyARN, roleName string) error {
	return nil
}

func (s *IAMService) ListRolePolicies(ctx context.Context, roleName string) ([]*awsoutputs.PolicyOutput, error) {
	return []*awsoutputs.PolicyOutput{}, nil
}

func (s *IAMService) AttachPolicyToGroup(ctx context.Context, policyARN, groupName string) error {
	return nil
}

func (s *IAMService) DetachPolicyFromGroup(ctx context.Context, policyARN, groupName string) error {
	return nil
}

func (s *IAMService) ListGroupPolicies(ctx context.Context, groupName string) ([]*awsoutputs.PolicyOutput, error) {
	return []*awsoutputs.PolicyOutput{}, nil
}

// User-Group operations

func (s *IAMService) AddUserToGroup(ctx context.Context, userName, groupName string) error {
	return nil
}

func (s *IAMService) RemoveUserFromGroup(ctx context.Context, userName, groupName string) error {
	return nil
}

func (s *IAMService) ListGroupUsers(ctx context.Context, groupName string) ([]*awsoutputs.UserOutput, error) {
	return []*awsoutputs.UserOutput{}, nil
}

// Inline Policy operations for Users

func (s *IAMService) PutUserInlinePolicy(ctx context.Context, userName string, policy *awsiam.InlinePolicy) error {
	return nil
}

func (s *IAMService) GetUserInlinePolicy(ctx context.Context, userName, policyName string) (*awsiam.InlinePolicy, error) {
	return &awsiam.InlinePolicy{
		Name:   policyName,
		Policy: `{"Version":"2012-10-17","Statement":[]}`,
	}, nil
}

func (s *IAMService) DeleteUserInlinePolicy(ctx context.Context, userName, policyName string) error {
	return nil
}

func (s *IAMService) ListUserInlinePolicies(ctx context.Context, userName string) ([]string, error) {
	return []string{}, nil
}

// Inline Policy operations for Roles

func (s *IAMService) PutRoleInlinePolicy(ctx context.Context, roleName string, policy *awsiam.InlinePolicy) error {
	return nil
}

func (s *IAMService) GetRoleInlinePolicy(ctx context.Context, roleName, policyName string) (*awsiam.InlinePolicy, error) {
	return &awsiam.InlinePolicy{
		Name:   policyName,
		Policy: `{"Version":"2012-10-17","Statement":[]}`,
	}, nil
}

func (s *IAMService) DeleteRoleInlinePolicy(ctx context.Context, roleName, policyName string) error {
	return nil
}

func (s *IAMService) ListRoleInlinePolicies(ctx context.Context, roleName string) ([]string, error) {
	return []string{}, nil
}

// Inline Policy operations for Groups

func (s *IAMService) PutGroupInlinePolicy(ctx context.Context, groupName string, policy *awsiam.InlinePolicy) error {
	return nil
}

func (s *IAMService) GetGroupInlinePolicy(ctx context.Context, groupName, policyName string) (*awsiam.InlinePolicy, error) {
	return &awsiam.InlinePolicy{
		Name:   policyName,
		Policy: `{"Version":"2012-10-17","Statement":[]}`,
	}, nil
}

func (s *IAMService) DeleteGroupInlinePolicy(ctx context.Context, groupName, policyName string) error {
	return nil
}

func (s *IAMService) ListGroupInlinePolicies(ctx context.Context, groupName string) ([]string, error) {
	return []string{}, nil
}

// Instance Profile operations

func (s *IAMService) CreateInstanceProfile(ctx context.Context, profile *awsiam.InstanceProfile) (*awsoutputs.InstanceProfileOutput, error) {
	if profile == nil {
		return nil, fmt.Errorf("instance profile is nil")
	}

	path := "/"
	if profile.Path != nil {
		path = *profile.Path
	}
	name := profile.Name
	if name == "" && profile.NamePrefix != nil {
		name = *profile.NamePrefix + "-12345"
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:instance-profile%s%s", path, name)
	if path == "" {
		arn = fmt.Sprintf("arn:aws:iam::123456789012:instance-profile/%s", name)
	}

	return &awsoutputs.InstanceProfileOutput{
		ARN:        arn,
		ID:         name,
		Name:       name,
		Path:       path,
		CreateDate: services.GetFixedTimestamp(),
		Tags:       profile.Tags,
		Roles:      []*awsoutputs.RoleOutput{},
	}, nil
}

func (s *IAMService) GetInstanceProfile(ctx context.Context, name string) (*awsoutputs.InstanceProfileOutput, error) {
	arn := fmt.Sprintf("arn:aws:iam::123456789012:instance-profile/%s", name)

	return &awsoutputs.InstanceProfileOutput{
		ARN:        arn,
		ID:         name,
		Name:       name,
		Path:       "/",
		CreateDate: services.GetFixedTimestamp(),
		Tags:       []configs.Tag{},
		Roles:      []*awsoutputs.RoleOutput{},
	}, nil
}

func (s *IAMService) UpdateInstanceProfile(ctx context.Context, name string, profile *awsiam.InstanceProfile) (*awsoutputs.InstanceProfileOutput, error) {
	return s.CreateInstanceProfile(ctx, profile)
}

func (s *IAMService) DeleteInstanceProfile(ctx context.Context, name string) error {
	return nil
}

func (s *IAMService) ListInstanceProfiles(ctx context.Context, pathPrefix *string) ([]*awsoutputs.InstanceProfileOutput, error) {
	path := "/"
	if pathPrefix != nil {
		path = *pathPrefix
	}
	arn := fmt.Sprintf("arn:aws:iam::123456789012:instance-profile%stest-profile", path)
	if path == "" {
		arn = "arn:aws:iam::123456789012:instance-profile/test-profile"
	}

	return []*awsoutputs.InstanceProfileOutput{
		{
			ARN:        arn,
			ID:         "test-profile",
			Name:       "test-profile",
			Path:       path,
			CreateDate: services.GetFixedTimestamp(),
			Tags:       []configs.Tag{},
			Roles:      []*awsoutputs.RoleOutput{},
		},
	}, nil
}

func (s *IAMService) AddRoleToInstanceProfile(ctx context.Context, profileName, roleName string) error {
	return nil
}

func (s *IAMService) RemoveRoleFromInstanceProfile(ctx context.Context, profileName, roleName string) error {
	return nil
}

func (s *IAMService) GetInstanceProfileRoles(ctx context.Context, profileName string) ([]*awsoutputs.RoleOutput, error) {
	return []*awsoutputs.RoleOutput{}, nil
}
