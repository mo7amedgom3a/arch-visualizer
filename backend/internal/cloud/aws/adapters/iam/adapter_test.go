package iam

import (
	"context"
	"errors"
	"testing"
	"time"

	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/iam"
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
)

// mockAWSIAMService is a mock implementation of AWSIAMService for testing
type mockAWSIAMService struct {
	policies            map[string]*awsoutputs.PolicyOutput
	roles               map[string]*awsoutputs.RoleOutput
	users               map[string]*awsoutputs.UserOutput
	groups              map[string]*awsoutputs.GroupOutput
	instanceProfiles    map[string]*awsoutputs.InstanceProfileOutput
	profileRoles        map[string][]string                        // profileName -> roleNames
	userPolicies        map[string][]string                        // userName -> policyARNs
	rolePolicies        map[string][]string                        // roleName -> policyARNs
	groupPolicies       map[string][]string                        // groupName -> policyARNs
	groupUsers          map[string][]string                        // groupName -> userNames
	userInlinePolicies  map[string]map[string]*awsiam.InlinePolicy // userName -> policyName -> policy
	roleInlinePolicies  map[string]map[string]*awsiam.InlinePolicy // roleName -> policyName -> policy
	groupInlinePolicies map[string]map[string]*awsiam.InlinePolicy // groupName -> policyName -> policy
	createError         error
	getError            error
}

// Ensure mockAWSIAMService implements AWSIAMService
var _ awsservice.AWSIAMService = (*mockAWSIAMService)(nil)

func newMockAWSIAMService() *mockAWSIAMService {
	return &mockAWSIAMService{
		policies:            make(map[string]*awsoutputs.PolicyOutput),
		roles:               make(map[string]*awsoutputs.RoleOutput),
		users:               make(map[string]*awsoutputs.UserOutput),
		groups:              make(map[string]*awsoutputs.GroupOutput),
		instanceProfiles:    make(map[string]*awsoutputs.InstanceProfileOutput),
		profileRoles:        make(map[string][]string),
		userPolicies:        make(map[string][]string),
		rolePolicies:        make(map[string][]string),
		groupPolicies:       make(map[string][]string),
		groupUsers:          make(map[string][]string),
		userInlinePolicies:  make(map[string]map[string]*awsiam.InlinePolicy),
		roleInlinePolicies:  make(map[string]map[string]*awsiam.InlinePolicy),
		groupInlinePolicies: make(map[string]map[string]*awsiam.InlinePolicy),
	}
}

// Policy Operations

func (m *mockAWSIAMService) CreatePolicy(ctx context.Context, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	arn := "arn:aws:iam::123456789012:policy/" + policy.Name
	output := &awsoutputs.PolicyOutput{
		ARN:            arn,
		ID:             arn,
		Name:           policy.Name,
		Description:    policy.Description,
		Path:           getPathOrDefault(policy.Path),
		PolicyDocument: policy.PolicyDocument,
		CreateDate:     time.Now(),
		UpdateDate:     time.Now(),
		Tags:           policy.Tags,
	}
	m.policies[arn] = output
	return output, nil
}

func (m *mockAWSIAMService) GetPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	policy, ok := m.policies[arn]
	if !ok {
		return nil, errors.New("policy not found")
	}
	return policy, nil
}

func (m *mockAWSIAMService) UpdatePolicy(ctx context.Context, arn string, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error) {
	existing, ok := m.policies[arn]
	if !ok {
		return nil, errors.New("policy not found")
	}
	existing.Name = policy.Name
	existing.Description = policy.Description
	existing.PolicyDocument = policy.PolicyDocument
	existing.UpdateDate = time.Now()
	existing.Tags = policy.Tags
	return existing, nil
}

func (m *mockAWSIAMService) DeletePolicy(ctx context.Context, arn string) error {
	delete(m.policies, arn)
	return nil
}

func (m *mockAWSIAMService) ListPolicies(ctx context.Context, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error) {
	var results []*awsoutputs.PolicyOutput
	for _, policy := range m.policies {
		if pathPrefix == nil || *pathPrefix == "" || policy.Path == *pathPrefix {
			results = append(results, policy)
		}
	}
	return results, nil
}

func (m *mockAWSIAMService) ListPoliciesBetweenServices(ctx context.Context, sourceService, destinationService string) ([]*awsoutputs.PolicyOutput, error) {
	// Simple mock implementation that returns nothing or could be customized
	return []*awsoutputs.PolicyOutput{}, nil
}

// Role Operations

func (m *mockAWSIAMService) CreateRole(ctx context.Context, role *awsiam.Role) (*awsoutputs.RoleOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	arn := "arn:aws:iam::123456789012:role/" + role.Name
	output := &awsoutputs.RoleOutput{
		ARN:                 arn,
		ID:                  role.Name,
		Name:                role.Name,
		Description:         role.Description,
		Path:                getPathOrDefault(role.Path),
		AssumeRolePolicy:    role.AssumeRolePolicy,
		PermissionsBoundary: role.PermissionsBoundary,
		CreateDate:          time.Now(),
		Tags:                role.Tags,
	}
	m.roles[role.Name] = output
	return output, nil
}

func (m *mockAWSIAMService) GetRole(ctx context.Context, name string) (*awsoutputs.RoleOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	role, ok := m.roles[name]
	if !ok {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func (m *mockAWSIAMService) UpdateRole(ctx context.Context, name string, role *awsiam.Role) (*awsoutputs.RoleOutput, error) {
	existing, ok := m.roles[name]
	if !ok {
		return nil, errors.New("role not found")
	}
	existing.Description = role.Description
	existing.AssumeRolePolicy = role.AssumeRolePolicy
	existing.PermissionsBoundary = role.PermissionsBoundary
	existing.Tags = role.Tags
	return existing, nil
}

func (m *mockAWSIAMService) DeleteRole(ctx context.Context, name string) error {
	delete(m.roles, name)
	return nil
}

func (m *mockAWSIAMService) ListRoles(ctx context.Context, pathPrefix *string) ([]*awsoutputs.RoleOutput, error) {
	var results []*awsoutputs.RoleOutput
	for _, role := range m.roles {
		if pathPrefix == nil || *pathPrefix == "" || role.Path == *pathPrefix {
			results = append(results, role)
		}
	}
	return results, nil
}

// User Operations

func (m *mockAWSIAMService) CreateUser(ctx context.Context, user *awsiam.User) (*awsoutputs.UserOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	arn := "arn:aws:iam::123456789012:user/" + user.Name
	output := &awsoutputs.UserOutput{
		ARN:                 arn,
		ID:                  user.Name,
		Name:                user.Name,
		Path:                getPathOrDefault(user.Path),
		PermissionsBoundary: user.PermissionsBoundary,
		CreateDate:          time.Now(),
		Tags:                user.Tags,
	}
	m.users[user.Name] = output
	return output, nil
}

func (m *mockAWSIAMService) GetUser(ctx context.Context, name string) (*awsoutputs.UserOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	user, ok := m.users[name]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockAWSIAMService) UpdateUser(ctx context.Context, name string, user *awsiam.User) (*awsoutputs.UserOutput, error) {
	existing, ok := m.users[name]
	if !ok {
		return nil, errors.New("user not found")
	}
	existing.Path = getPathOrDefault(user.Path)
	existing.PermissionsBoundary = user.PermissionsBoundary
	existing.Tags = user.Tags
	return existing, nil
}

func (m *mockAWSIAMService) DeleteUser(ctx context.Context, name string) error {
	delete(m.users, name)
	return nil
}

func (m *mockAWSIAMService) ListUsers(ctx context.Context, pathPrefix *string) ([]*awsoutputs.UserOutput, error) {
	var results []*awsoutputs.UserOutput
	for _, user := range m.users {
		if pathPrefix == nil || *pathPrefix == "" || user.Path == *pathPrefix {
			results = append(results, user)
		}
	}
	return results, nil
}

// Group Operations

func (m *mockAWSIAMService) CreateGroup(ctx context.Context, group *awsiam.Group) (*awsoutputs.GroupOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	arn := "arn:aws:iam::123456789012:group/" + group.Name
	output := &awsoutputs.GroupOutput{
		ARN:        arn,
		ID:         group.Name,
		Name:       group.Name,
		Path:       getPathOrDefault(group.Path),
		CreateDate: time.Now(),
		Tags:       group.Tags,
	}
	m.groups[group.Name] = output
	return output, nil
}

func (m *mockAWSIAMService) GetGroup(ctx context.Context, name string) (*awsoutputs.GroupOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	group, ok := m.groups[name]
	if !ok {
		return nil, errors.New("group not found")
	}
	return group, nil
}

func (m *mockAWSIAMService) UpdateGroup(ctx context.Context, name string, group *awsiam.Group) (*awsoutputs.GroupOutput, error) {
	existing, ok := m.groups[name]
	if !ok {
		return nil, errors.New("group not found")
	}
	existing.Path = getPathOrDefault(group.Path)
	existing.Tags = group.Tags
	return existing, nil
}

func (m *mockAWSIAMService) DeleteGroup(ctx context.Context, name string) error {
	delete(m.groups, name)
	return nil
}

func (m *mockAWSIAMService) ListGroups(ctx context.Context, pathPrefix *string) ([]*awsoutputs.GroupOutput, error) {
	var results []*awsoutputs.GroupOutput
	for _, group := range m.groups {
		if pathPrefix == nil || *pathPrefix == "" || group.Path == *pathPrefix {
			results = append(results, group)
		}
	}
	return results, nil
}

// Policy Attachment Operations

func (m *mockAWSIAMService) AttachPolicyToUser(ctx context.Context, policyARN, userName string) error {
	if _, ok := m.users[userName]; !ok {
		return errors.New("user not found")
	}
	if _, ok := m.policies[policyARN]; !ok {
		return errors.New("policy not found")
	}
	m.userPolicies[userName] = append(m.userPolicies[userName], policyARN)
	return nil
}

func (m *mockAWSIAMService) DetachPolicyFromUser(ctx context.Context, policyARN, userName string) error {
	policies := m.userPolicies[userName]
	for i, arn := range policies {
		if arn == policyARN {
			m.userPolicies[userName] = append(policies[:i], policies[i+1:]...)
			return nil
		}
	}
	return errors.New("policy not attached to user")
}

func (m *mockAWSIAMService) ListUserPolicies(ctx context.Context, userName string) ([]*awsoutputs.PolicyOutput, error) {
	policyARNs := m.userPolicies[userName]
	var results []*awsoutputs.PolicyOutput
	for _, arn := range policyARNs {
		if policy, ok := m.policies[arn]; ok {
			results = append(results, policy)
		}
	}
	return results, nil
}

func (m *mockAWSIAMService) AttachPolicyToRole(ctx context.Context, policyARN, roleName string) error {
	if _, ok := m.roles[roleName]; !ok {
		return errors.New("role not found")
	}
	if _, ok := m.policies[policyARN]; !ok {
		return errors.New("policy not found")
	}
	m.rolePolicies[roleName] = append(m.rolePolicies[roleName], policyARN)
	return nil
}

func (m *mockAWSIAMService) DetachPolicyFromRole(ctx context.Context, policyARN, roleName string) error {
	policies := m.rolePolicies[roleName]
	for i, arn := range policies {
		if arn == policyARN {
			m.rolePolicies[roleName] = append(policies[:i], policies[i+1:]...)
			return nil
		}
	}
	return errors.New("policy not attached to role")
}

func (m *mockAWSIAMService) ListRolePolicies(ctx context.Context, roleName string) ([]*awsoutputs.PolicyOutput, error) {
	policyARNs := m.rolePolicies[roleName]
	var results []*awsoutputs.PolicyOutput
	for _, arn := range policyARNs {
		if policy, ok := m.policies[arn]; ok {
			results = append(results, policy)
		}
	}
	return results, nil
}

func (m *mockAWSIAMService) AttachPolicyToGroup(ctx context.Context, policyARN, groupName string) error {
	if _, ok := m.groups[groupName]; !ok {
		return errors.New("group not found")
	}
	if _, ok := m.policies[policyARN]; !ok {
		return errors.New("policy not found")
	}
	m.groupPolicies[groupName] = append(m.groupPolicies[groupName], policyARN)
	return nil
}

func (m *mockAWSIAMService) DetachPolicyFromGroup(ctx context.Context, policyARN, groupName string) error {
	policies := m.groupPolicies[groupName]
	for i, arn := range policies {
		if arn == policyARN {
			m.groupPolicies[groupName] = append(policies[:i], policies[i+1:]...)
			return nil
		}
	}
	return errors.New("policy not attached to group")
}

func (m *mockAWSIAMService) ListGroupPolicies(ctx context.Context, groupName string) ([]*awsoutputs.PolicyOutput, error) {
	policyARNs := m.groupPolicies[groupName]
	var results []*awsoutputs.PolicyOutput
	for _, arn := range policyARNs {
		if policy, ok := m.policies[arn]; ok {
			results = append(results, policy)
		}
	}
	return results, nil
}

// User-Group Operations

func (m *mockAWSIAMService) AddUserToGroup(ctx context.Context, userName, groupName string) error {
	if _, ok := m.users[userName]; !ok {
		return errors.New("user not found")
	}
	if _, ok := m.groups[groupName]; !ok {
		return errors.New("group not found")
	}
	m.groupUsers[groupName] = append(m.groupUsers[groupName], userName)
	return nil
}

func (m *mockAWSIAMService) RemoveUserFromGroup(ctx context.Context, userName, groupName string) error {
	users := m.groupUsers[groupName]
	for i, name := range users {
		if name == userName {
			m.groupUsers[groupName] = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not in group")
}

func (m *mockAWSIAMService) ListGroupUsers(ctx context.Context, groupName string) ([]*awsoutputs.UserOutput, error) {
	userNames := m.groupUsers[groupName]
	var results []*awsoutputs.UserOutput
	for _, name := range userNames {
		if user, ok := m.users[name]; ok {
			results = append(results, user)
		}
	}
	return results, nil
}

// Inline Policy Operations

func (m *mockAWSIAMService) PutUserInlinePolicy(ctx context.Context, userName string, policy *awsiam.InlinePolicy) error {
	if _, ok := m.users[userName]; !ok {
		return errors.New("user not found")
	}
	if m.userInlinePolicies[userName] == nil {
		m.userInlinePolicies[userName] = make(map[string]*awsiam.InlinePolicy)
	}
	m.userInlinePolicies[userName][policy.Name] = policy
	return nil
}

func (m *mockAWSIAMService) GetUserInlinePolicy(ctx context.Context, userName, policyName string) (*awsiam.InlinePolicy, error) {
	if policies, ok := m.userInlinePolicies[userName]; ok {
		if policy, ok := policies[policyName]; ok {
			return policy, nil
		}
	}
	return nil, errors.New("inline policy not found")
}

func (m *mockAWSIAMService) DeleteUserInlinePolicy(ctx context.Context, userName, policyName string) error {
	if policies, ok := m.userInlinePolicies[userName]; ok {
		delete(policies, policyName)
		return nil
	}
	return errors.New("inline policy not found")
}

func (m *mockAWSIAMService) ListUserInlinePolicies(ctx context.Context, userName string) ([]string, error) {
	var names []string
	if policies, ok := m.userInlinePolicies[userName]; ok {
		for name := range policies {
			names = append(names, name)
		}
	}
	return names, nil
}

func (m *mockAWSIAMService) PutRoleInlinePolicy(ctx context.Context, roleName string, policy *awsiam.InlinePolicy) error {
	if _, ok := m.roles[roleName]; !ok {
		return errors.New("role not found")
	}
	if m.roleInlinePolicies[roleName] == nil {
		m.roleInlinePolicies[roleName] = make(map[string]*awsiam.InlinePolicy)
	}
	m.roleInlinePolicies[roleName][policy.Name] = policy
	return nil
}

func (m *mockAWSIAMService) GetRoleInlinePolicy(ctx context.Context, roleName, policyName string) (*awsiam.InlinePolicy, error) {
	if policies, ok := m.roleInlinePolicies[roleName]; ok {
		if policy, ok := policies[policyName]; ok {
			return policy, nil
		}
	}
	return nil, errors.New("inline policy not found")
}

func (m *mockAWSIAMService) DeleteRoleInlinePolicy(ctx context.Context, roleName, policyName string) error {
	if policies, ok := m.roleInlinePolicies[roleName]; ok {
		delete(policies, policyName)
		return nil
	}
	return errors.New("inline policy not found")
}

func (m *mockAWSIAMService) ListRoleInlinePolicies(ctx context.Context, roleName string) ([]string, error) {
	var names []string
	if policies, ok := m.roleInlinePolicies[roleName]; ok {
		for name := range policies {
			names = append(names, name)
		}
	}
	return names, nil
}

func (m *mockAWSIAMService) PutGroupInlinePolicy(ctx context.Context, groupName string, policy *awsiam.InlinePolicy) error {
	if _, ok := m.groups[groupName]; !ok {
		return errors.New("group not found")
	}
	if m.groupInlinePolicies[groupName] == nil {
		m.groupInlinePolicies[groupName] = make(map[string]*awsiam.InlinePolicy)
	}
	m.groupInlinePolicies[groupName][policy.Name] = policy
	return nil
}

func (m *mockAWSIAMService) GetGroupInlinePolicy(ctx context.Context, groupName, policyName string) (*awsiam.InlinePolicy, error) {
	if policies, ok := m.groupInlinePolicies[groupName]; ok {
		if policy, ok := policies[policyName]; ok {
			return policy, nil
		}
	}
	return nil, errors.New("inline policy not found")
}

func (m *mockAWSIAMService) DeleteGroupInlinePolicy(ctx context.Context, groupName, policyName string) error {
	if policies, ok := m.groupInlinePolicies[groupName]; ok {
		delete(policies, policyName)
		return nil
	}
	return errors.New("inline policy not found")
}

func (m *mockAWSIAMService) ListGroupInlinePolicies(ctx context.Context, groupName string) ([]string, error) {
	var names []string
	if policies, ok := m.groupInlinePolicies[groupName]; ok {
		for name := range policies {
			names = append(names, name)
		}
	}
	return names, nil
}

// AWS Managed Policy Operations

func (m *mockAWSIAMService) ListAWSManagedPolicies(ctx context.Context, scope *string, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error) {
	// Return some mock AWS managed policies
	return []*awsoutputs.PolicyOutput{
		{
			ARN:            "arn:aws:iam::aws:policy/ReadOnlyAccess",
			ID:             "arn:aws:iam::aws:policy/ReadOnlyAccess",
			Name:           "ReadOnlyAccess",
			Path:           "/",
			PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`,
			CreateDate:     time.Now(),
			UpdateDate:     time.Now(),
			IsAttachable:   true,
			IsAWSManaged:   true,
		},
		{
			ARN:            "arn:aws:iam::aws:policy/PowerUserAccess",
			ID:             "arn:aws:iam::aws:policy/PowerUserAccess",
			Name:           "PowerUserAccess",
			Path:           "/",
			PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`,
			CreateDate:     time.Now(),
			UpdateDate:     time.Now(),
			IsAttachable:   true,
			IsAWSManaged:   true,
		},
	}, nil
}

func (m *mockAWSIAMService) GetAWSManagedPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error) {
	// Return a mock AWS managed policy
	if arn == "arn:aws:iam::aws:policy/ReadOnlyAccess" {
		return &awsoutputs.PolicyOutput{
			ARN:            arn,
			ID:             arn,
			Name:           "ReadOnlyAccess",
			Path:           "/",
			PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`,
			CreateDate:     time.Now(),
			UpdateDate:     time.Now(),
			IsAttachable:   true,
			IsAWSManaged:   true,
		}, nil
	}
	return nil, errors.New("AWS managed policy not found")
}

// Instance Profile Operations

func (m *mockAWSIAMService) CreateInstanceProfile(ctx context.Context, profile *awsiam.InstanceProfile) (*awsoutputs.InstanceProfileOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	name := profile.Name
	if name == "" && profile.NamePrefix != nil {
		name = *profile.NamePrefix + "-12345"
	}
	arn := "arn:aws:iam::123456789012:instance-profile/" + name
	output := &awsoutputs.InstanceProfileOutput{
		ARN:        arn,
		ID:         name,
		Name:       name,
		Path:       getPathOrDefault(profile.Path),
		CreateDate: time.Now(),
		Tags:       profile.Tags,
		Roles:      make([]*awsoutputs.RoleOutput, 0),
	}
	if profile.Role != nil && *profile.Role != "" {
		// Add role if provided
		if role, ok := m.roles[*profile.Role]; ok {
			output.Roles = []*awsoutputs.RoleOutput{role}
			m.profileRoles[name] = []string{*profile.Role}
		}
	}
	m.instanceProfiles[name] = output
	return output, nil
}

func (m *mockAWSIAMService) GetInstanceProfile(ctx context.Context, name string) (*awsoutputs.InstanceProfileOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	profile, ok := m.instanceProfiles[name]
	if !ok {
		return nil, errors.New("instance profile not found")
	}
	return profile, nil
}

func (m *mockAWSIAMService) UpdateInstanceProfile(ctx context.Context, name string, profile *awsiam.InstanceProfile) (*awsoutputs.InstanceProfileOutput, error) {
	existing, ok := m.instanceProfiles[name]
	if !ok {
		return nil, errors.New("instance profile not found")
	}
	existing.Tags = profile.Tags
	return existing, nil
}

func (m *mockAWSIAMService) DeleteInstanceProfile(ctx context.Context, name string) error {
	delete(m.instanceProfiles, name)
	delete(m.profileRoles, name)
	return nil
}

func (m *mockAWSIAMService) ListInstanceProfiles(ctx context.Context, pathPrefix *string) ([]*awsoutputs.InstanceProfileOutput, error) {
	var results []*awsoutputs.InstanceProfileOutput
	for _, profile := range m.instanceProfiles {
		if pathPrefix == nil || *pathPrefix == "" || profile.Path == *pathPrefix {
			results = append(results, profile)
		}
	}
	return results, nil
}

func (m *mockAWSIAMService) AddRoleToInstanceProfile(ctx context.Context, profileName, roleName string) error {
	if _, ok := m.instanceProfiles[profileName]; !ok {
		return errors.New("instance profile not found")
	}
	if _, ok := m.roles[roleName]; !ok {
		return errors.New("role not found")
	}
	// Add role to profile
	role := m.roles[roleName]
	profile := m.instanceProfiles[profileName]
	profile.Roles = append(profile.Roles, role)
	m.profileRoles[profileName] = append(m.profileRoles[profileName], roleName)
	return nil
}

func (m *mockAWSIAMService) RemoveRoleFromInstanceProfile(ctx context.Context, profileName, roleName string) error {
	if _, ok := m.instanceProfiles[profileName]; !ok {
		return errors.New("instance profile not found")
	}
	roles := m.profileRoles[profileName]
	for i, name := range roles {
		if name == roleName {
			m.profileRoles[profileName] = append(roles[:i], roles[i+1:]...)
			// Remove from profile.Roles slice
			profile := m.instanceProfiles[profileName]
			for j, role := range profile.Roles {
				if role != nil && role.Name == roleName {
					profile.Roles = append(profile.Roles[:j], profile.Roles[j+1:]...)
					break
				}
			}
			return nil
		}
	}
	return errors.New("role not attached to instance profile")
}

func (m *mockAWSIAMService) GetInstanceProfileRoles(ctx context.Context, profileName string) ([]*awsoutputs.RoleOutput, error) {
	profile, ok := m.instanceProfiles[profileName]
	if !ok {
		return nil, errors.New("instance profile not found")
	}
	return profile.Roles, nil
}

// Helper function
func getPathOrDefault(path *string) string {
	if path != nil && *path != "" {
		return *path
	}
	return "/"
}

// Tests

func TestAWSIAMAdapter_CreatePolicy(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	domainPolicy := &domainiam.Policy{
		Name:           "test-policy",
		PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`,
	}

	ctx := context.Background()
	createdPolicy, err := adapter.CreatePolicy(ctx, domainPolicy)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdPolicy == nil {
		t.Fatal("Expected created policy, got nil")
	}

	if createdPolicy.Name != domainPolicy.Name {
		t.Errorf("Expected name %s, got %s", domainPolicy.Name, createdPolicy.Name)
	}

	if createdPolicy.ARN == nil || *createdPolicy.ARN == "" {
		t.Error("Expected policy ARN to be populated")
	}
}

func TestAWSIAMAdapter_CreateRole(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	domainRole := &domainiam.Role{
		Name:             "test-role",
		AssumeRolePolicy: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}`,
	}

	ctx := context.Background()
	createdRole, err := adapter.CreateRole(ctx, domainRole)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdRole == nil {
		t.Fatal("Expected created role, got nil")
	}

	if createdRole.Name != domainRole.Name {
		t.Errorf("Expected name %s, got %s", domainRole.Name, createdRole.Name)
	}

	if createdRole.ARN == nil || *createdRole.ARN == "" {
		t.Error("Expected role ARN to be populated")
	}
}

func TestAWSIAMAdapter_CreateUser(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	domainUser := &domainiam.User{
		Name: "test-user",
	}

	ctx := context.Background()
	createdUser, err := adapter.CreateUser(ctx, domainUser)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdUser == nil {
		t.Fatal("Expected created user, got nil")
	}

	if createdUser.Name != domainUser.Name {
		t.Errorf("Expected name %s, got %s", domainUser.Name, createdUser.Name)
	}

	if createdUser.ARN == nil || *createdUser.ARN == "" {
		t.Error("Expected user ARN to be populated")
	}
}

func TestAWSIAMAdapter_CreateGroup(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	domainGroup := &domainiam.Group{
		Name: "test-group",
	}

	ctx := context.Background()
	createdGroup, err := adapter.CreateGroup(ctx, domainGroup)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdGroup == nil {
		t.Fatal("Expected created group, got nil")
	}

	if createdGroup.Name != domainGroup.Name {
		t.Errorf("Expected name %s, got %s", domainGroup.Name, createdGroup.Name)
	}

	if createdGroup.ARN == nil || *createdGroup.ARN == "" {
		t.Error("Expected group ARN to be populated")
	}
}

func TestAWSIAMAdapter_AttachPolicyToUser(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	ctx := context.Background()

	// Create policy and user first
	policy := &domainiam.Policy{
		Name:           "test-policy",
		PolicyDocument: `{"Version":"2012-10-17","Statement":[]}`,
	}
	createdPolicy, _ := adapter.CreatePolicy(ctx, policy)

	user := &domainiam.User{Name: "test-user"}
	adapter.CreateUser(ctx, user)

	// Attach policy
	err := adapter.AttachPolicyToUser(ctx, *createdPolicy.ARN, user.Name)
	if err != nil {
		t.Fatalf("Expected no error attaching policy, got: %v", err)
	}

	// Verify attachment
	policies, err := adapter.ListUserPolicies(ctx, user.Name)
	if err != nil {
		t.Fatalf("Expected no error listing policies, got: %v", err)
	}

	if len(policies) != 1 {
		t.Errorf("Expected 1 policy attached, got %d", len(policies))
	}
}

func TestAWSIAMAdapter_AddUserToGroup(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	ctx := context.Background()

	// Create user and group
	user := &domainiam.User{Name: "test-user"}
	adapter.CreateUser(ctx, user)

	group := &domainiam.Group{Name: "test-group"}
	adapter.CreateGroup(ctx, group)

	// Add user to group
	err := adapter.AddUserToGroup(ctx, user.Name, group.Name)
	if err != nil {
		t.Fatalf("Expected no error adding user to group, got: %v", err)
	}

	// Verify membership
	users, err := adapter.ListGroupUsers(ctx, group.Name)
	if err != nil {
		t.Fatalf("Expected no error listing group users, got: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user in group, got %d", len(users))
	}
}

func TestAWSIAMAdapter_ValidationError(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	ctx := context.Background()

	// Try to create policy without name
	invalidPolicy := &domainiam.Policy{
		PolicyDocument: `{"Version":"2012-10-17","Statement":[]}`,
	}

	_, err := adapter.CreatePolicy(ctx, invalidPolicy)
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	if !contains(err.Error(), "validation failed") {
		t.Errorf("Expected validation error, got: %v", err)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestAWSIAMAdapter_PutRoleInlinePolicy(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	ctx := context.Background()

	// Create role first
	role := &domainiam.Role{
		Name:             "test-role",
		AssumeRolePolicy: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}`,
	}
	adapter.CreateRole(ctx, role)

	// Create inline policy
	inlinePolicy := &domainiam.InlinePolicy{
		Name:   "test-inline-policy",
		Policy: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`,
	}

	err := adapter.PutRoleInlinePolicy(ctx, role.Name, inlinePolicy)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify policy was created
	retrieved, err := adapter.GetRoleInlinePolicy(ctx, role.Name, inlinePolicy.Name)
	if err != nil {
		t.Fatalf("Expected no error retrieving policy, got: %v", err)
	}

	if retrieved.Name != inlinePolicy.Name {
		t.Errorf("Expected policy name %s, got %s", inlinePolicy.Name, retrieved.Name)
	}
}

func TestAWSIAMAdapter_ListAWSManagedPolicies(t *testing.T) {
	mockService := newMockAWSIAMService()
	adapter := NewAWSIAMAdapter(mockService)

	ctx := context.Background()
	policies, err := adapter.ListAWSManagedPolicies(ctx, nil, nil)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(policies) == 0 {
		t.Error("Expected at least one AWS managed policy")
	}

	// Verify policy type
	for _, policy := range policies {
		if policy.Type != domainiam.PolicyTypeAWSManaged {
			t.Errorf("Expected policy type %s, got %s", domainiam.PolicyTypeAWSManaged, policy.Type)
		}
	}
}
