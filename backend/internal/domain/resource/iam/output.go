package iam

import (
	"time"
)

// PolicyOutput represents the output data for an IAM policy after creation/update
type PolicyOutput struct {
	// Core identifiers
	ID            string
	ARN           *string
	Name          string
	Description   *string
	Path          *string

	// Configuration
	PolicyDocument string
	Tags          []PolicyTag
	Type          PolicyType
	IsAttachable  *bool

	// Output fields (cloud-generated)
	DefaultVersionID *string
	AttachmentCount  *int
	IsAWSManaged     *bool

	// Timestamps
	CreateDate *time.Time
	UpdateDate *time.Time
}

// RoleOutput represents the output data for an IAM role after creation/update
type RoleOutput struct {
	// Core identifiers
	ID                string
	ARN               *string
	Name              string
	Description       *string
	Path              *string
	UniqueID          *string // Cloud-generated unique identifier

	// Configuration
	AssumeRolePolicy  string
	PermissionsBoundary *string
	Tags              []PolicyTag

	// Output fields (cloud-generated)
	MaxSessionDuration *int

	// Timestamps
	CreateDate *time.Time
}

// UserOutput represents the output data for an IAM user after creation/update
type UserOutput struct {
	// Core identifiers
	ID                string
	ARN               *string
	Name              string
	Path              *string
	UniqueID          *string // Cloud-generated unique identifier

	// Configuration
	PermissionsBoundary *string
	ForceDestroy      *bool
	Tags              []PolicyTag

	// Timestamps
	CreateDate *time.Time
}

// GroupOutput represents the output data for an IAM group after creation/update
type GroupOutput struct {
	// Core identifiers
	ID       string
	ARN      *string
	Name     string
	Path     *string
	UniqueID *string // Cloud-generated unique identifier

	// Configuration
	Tags []PolicyTag

	// Timestamps
	CreateDate *time.Time
}

// InstanceProfileOutput represents the output data for an IAM instance profile after creation/update
type InstanceProfileOutput struct {
	// Core identifiers
	ID       string
	ARN      *string
	Name     string
	Path     *string
	UniqueID *string // Cloud-generated unique identifier

	// Configuration
	RoleName *string
	Tags     []PolicyTag

	// Timestamps
	CreateDate *time.Time
}

// InlinePolicyOutput represents the output data for an inline policy
type InlinePolicyOutput struct {
	// Core identifiers
	Name     string
	PolicyDocument string

	// Timestamps
	CreateDate *time.Time
	UpdateDate *time.Time
}

// ToPolicyOutput converts a Policy domain model to PolicyOutput
func ToPolicyOutput(policy *Policy) *PolicyOutput {
	if policy == nil {
		return nil
	}
	return &PolicyOutput{
		ID:            policy.ID,
		ARN:           policy.ARN,
		Name:          policy.Name,
		Description:   policy.Description,
		Path:          policy.Path,
		PolicyDocument: policy.PolicyDocument,
		Tags:          policy.Tags,
		Type:          policy.Type,
		IsAttachable:  policy.IsAttachable,
	}
}

// ToRoleOutput converts a Role domain model to RoleOutput
func ToRoleOutput(role *Role) *RoleOutput {
	if role == nil {
		return nil
	}
	return &RoleOutput{
		ID:                role.ID,
		ARN:               role.ARN,
		Name:              role.Name,
		Description:       role.Description,
		Path:              role.Path,
		AssumeRolePolicy:  role.AssumeRolePolicy,
		PermissionsBoundary: role.PermissionsBoundary,
		Tags:              role.Tags,
	}
}

// ToUserOutput converts a User domain model to UserOutput
func ToUserOutput(user *User) *UserOutput {
	if user == nil {
		return nil
	}
	return &UserOutput{
		ID:                user.ID,
		ARN:               user.ARN,
		Name:              user.Name,
		Path:              user.Path,
		PermissionsBoundary: user.PermissionsBoundary,
		ForceDestroy:      user.ForceDestroy,
		Tags:              user.Tags,
	}
}

// ToGroupOutput converts a Group domain model to GroupOutput
func ToGroupOutput(group *Group) *GroupOutput {
	if group == nil {
		return nil
	}
	return &GroupOutput{
		ID:   group.ID,
		ARN:  group.ARN,
		Name: group.Name,
		Path: group.Path,
		Tags: group.Tags,
	}
}

// ToInstanceProfileOutput converts an InstanceProfile domain model to InstanceProfileOutput
func ToInstanceProfileOutput(profile *InstanceProfile) *InstanceProfileOutput {
	if profile == nil {
		return nil
	}
	return &InstanceProfileOutput{
		ID:       profile.ID,
		ARN:      profile.ARN,
		Name:     profile.Name,
		Path:     profile.Path,
		RoleName: profile.RoleName,
		Tags:     profile.Tags,
	}
}
