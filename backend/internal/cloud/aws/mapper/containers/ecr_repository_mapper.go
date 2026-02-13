package containers

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// ECRRepositoryFromResource converts a generic domain resource to an ECR Repository model
func ECRRepositoryFromResource(res *resource.Resource) (*containers.ECRRepository, error) {
	if res.Type.Name != "ECRRepository" {
		return nil, fmt.Errorf("invalid resource type for ECR Repository mapper: %s", res.Type.Name)
	}

	repo := &containers.ECRRepository{
		Name: res.Name,
	}

	getString := func(key string) string {
		if val, ok := res.Metadata[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	getBool := func(key string) bool {
		if val, ok := res.Metadata[key]; ok {
			if b, ok := val.(bool); ok {
				return b
			}
		}
		return false
	}

	repo.ImageTagMutability = getString("image_tag_mutability")
	if repo.ImageTagMutability == "" {
		repo.ImageTagMutability = "MUTABLE"
	}

	repo.ScanOnPush = getBool("scan_on_push")
	repo.EncryptionType = getString("encryption_type")
	repo.KMSKey = getString("kms_key")
	repo.ForceDelete = getBool("force_delete")

	return repo, nil
}

// MapECRRepository maps an ECR Repository to a TerraformBlock
func MapECRRepository(repo *containers.ECRRepository) (*mapper.TerraformBlock, error) {
	if repo == nil {
		return nil, fmt.Errorf("ecr repository is nil")
	}

	attributes := make(map[string]mapper.TerraformValue)
	nestedBlocks := make(map[string][]mapper.NestedBlock)

	attributes["name"] = strVal(repo.Name)
	attributes["image_tag_mutability"] = strVal(repo.ImageTagMutability)

	if repo.ForceDelete {
		attributes["force_delete"] = boolVal(true)
	}

	// Image scanning configuration
	if repo.ScanOnPush {
		nestedBlocks["image_scanning_configuration"] = []mapper.NestedBlock{
			{
				Attributes: map[string]mapper.TerraformValue{
					"scan_on_push": boolVal(true),
				},
			},
		}
	}

	// Encryption configuration
	if repo.EncryptionType != "" {
		encAttrs := map[string]mapper.TerraformValue{
			"encryption_type": strVal(repo.EncryptionType),
		}
		if repo.KMSKey != "" {
			encAttrs["kms_key"] = strVal(repo.KMSKey)
		}
		nestedBlocks["encryption_configuration"] = []mapper.NestedBlock{
			{Attributes: encAttrs},
		}
	}

	return &mapper.TerraformBlock{
		Kind:         "resource",
		Labels:       []string{"aws_ecr_repository", repo.Name},
		Attributes:   attributes,
		NestedBlocks: nestedBlocks,
	}, nil
}
