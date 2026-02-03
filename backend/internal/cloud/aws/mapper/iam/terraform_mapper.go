package iam

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func init() {
	inv := inventory.GetDefaultInventory()
	inv.SetTerraformMapper("IAMPolicy", MapIAMPolicyToTerraform)
	inv.SetTerraformMapper("IAMRolePolicyAttachment", MapIAMRolePolicyAttachmentToTerraform)
	inv.SetTerraformMapper("IAMUser", MapIAMUserToTerraform)
	inv.SetTerraformMapper("IAMRole", MapIAMRoleToTerraform)
}

// MapIAMPolicyToTerraform resources to Terraform blocks
func MapIAMPolicyToTerraform(res *resource.Resource) ([]mapper.TerraformBlock, error) {
	name, ok := res.Metadata["name"].(string)
	if !ok || name == "" {
		name = res.Name
	}

	policyRaw, ok := res.Metadata["policy"]
	var policyDoc string
	if ok {
		policyDoc = fmt.Sprintf("%v", policyRaw)
	}

	block := mapper.TerraformBlock{
		Kind:   "resource",
		Labels: []string{"aws_iam_policy", name},
		Attributes: map[string]mapper.TerraformValue{
			"name":        strVal(name),
			"path":        strVal("/"),
			"description": strVal(fmt.Sprintf("%v", res.Metadata["description"])),
			"policy":      strVal(policyDoc),
		},
	}

	return []mapper.TerraformBlock{block}, nil
}

// MapIAMRolePolicyAttachmentToTerraform resources to Terraform blocks
func MapIAMRolePolicyAttachmentToTerraform(res *resource.Resource) ([]mapper.TerraformBlock, error) {
	name, ok := res.Metadata["name"].(string)
	if !ok || name == "" {
		name = res.Name
	}

	role := fmt.Sprintf("%v", res.Metadata["role"])
	policyARN := fmt.Sprintf("%v", res.Metadata["policy_arn"])

	block := mapper.TerraformBlock{
		Kind:   "resource",
		Labels: []string{"aws_iam_role_policy_attachment", name},
		Attributes: map[string]mapper.TerraformValue{
			"role":       strVal(role),
			"policy_arn": exprVal(policyARN),
		},
	}

	return []mapper.TerraformBlock{block}, nil
}

func strVal(s string) mapper.TerraformValue {
	return mapper.TerraformValue{String: &s}
}

func exprVal(s string) mapper.TerraformValue {
	e := mapper.TerraformExpr(s)
	return mapper.TerraformValue{Expr: &e}
}

// MapIAMUserToTerraform resources to Terraform blocks
func MapIAMUserToTerraform(res *resource.Resource) ([]mapper.TerraformBlock, error) {
	name, ok := res.Metadata["name"].(string)
	if !ok || name == "" {
		name = res.Name
	}

	block := mapper.TerraformBlock{
		Kind:   "resource",
		Labels: []string{"aws_iam_user", name},
		Attributes: map[string]mapper.TerraformValue{
			"name": strVal(name),
		},
	}
	if p, ok := res.Metadata["path"].(string); ok && p != "" {
		block.Attributes["path"] = strVal(p)
	}
	if p, ok := res.Metadata["permissions_boundary"].(string); ok && p != "" {
		block.Attributes["permissions_boundary"] = strVal(p)
	}
	if v, ok := res.Metadata["force_destroy"].(bool); ok && v {
		block.Attributes["force_destroy"] = mapper.TerraformValue{Bool: &v}
	}

	return []mapper.TerraformBlock{block}, nil
}

// MapIAMRoleToTerraform resources to Terraform blocks
func MapIAMRoleToTerraform(res *resource.Resource) ([]mapper.TerraformBlock, error) {
	name, ok := res.Metadata["name"].(string)
	if !ok || name == "" {
		name = res.Name
	}

	assumePolicy := "{}" // Default or Error?
	if p, ok := res.Metadata["assume_role_policy"].(string); ok && p != "" {
		assumePolicy = p
	}

	block := mapper.TerraformBlock{
		Kind:   "resource",
		Labels: []string{"aws_iam_role", name},
		Attributes: map[string]mapper.TerraformValue{
			"name":               strVal(name),
			"assume_role_policy": strVal(assumePolicy),
		},
	}
	if p, ok := res.Metadata["path"].(string); ok && p != "" {
		block.Attributes["path"] = strVal(p)
	}
	if desc, ok := res.Metadata["description"].(string); ok && desc != "" {
		block.Attributes["description"] = strVal(desc)
	}
	if p, ok := res.Metadata["permissions_boundary"].(string); ok && p != "" {
		block.Attributes["permissions_boundary"] = strVal(p)
	}

	return []mapper.TerraformBlock{block}, nil
}
