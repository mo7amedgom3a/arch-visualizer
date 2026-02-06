package database

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func init() {
	inv := inventory.GetDefaultInventory()

	// Register Terraform Mapper for RDS
	inv.SetTerraformMapper("RDS", func(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
		// Convert generic resource to AWS RDS Model
		rdsInstance, err := FromResource(res)
		if err != nil {
			return nil, err
		}

		block, err := MapRDSInstance(rdsInstance)
		if err != nil {
			return nil, err
		}
		return []tfmapper.TerraformBlock{*block}, nil
	})
}
