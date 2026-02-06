package database

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// FromResource converts a generic domain resource to an AWS RDS Instance model
func FromResource(res *resource.Resource) (*database.RDSInstance, error) {
	if res.Type.Name != "RDS" {
		return nil, fmt.Errorf("invalid resource type for RDS mapper: %s", res.Type.Name)
	}

	rds := &database.RDSInstance{
		Name: res.Name,
		// Tags: []configs.Tag{{Key: "Name", Value: res.Name}}, // Initial tags
	}
	// Initialize tags if needed or handle later

	// Helper to safely extract string from metadata
	getString := func(key string) string {
		if val, ok := res.Metadata[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	// Helper to safely extract int from metadata
	getInt := func(key string) int {
		if val, ok := res.Metadata[key]; ok {
			if i, ok := val.(int); ok {
				return i
			}
			if f, ok := val.(float64); ok {
				return int(f)
			}
		}
		return 0
	}

	// Helper to safely extract bool from metadata
	getBool := func(key string) bool {
		if val, ok := res.Metadata[key]; ok {
			if b, ok := val.(bool); ok {
				return b
			}
		}
		return false
	}

	rds.Engine = getString("engine")
	rds.EngineVersion = getString("engine_version")
	rds.InstanceClass = getString("instance_class")
	rds.AllocatedStorage = getInt("allocated_storage")
	rds.StorageType = getString("storage_type")
	rds.Username = getString("username")
	rds.Password = getString("password")
	rds.DBName = getString("db_name")
	rds.SubnetGroupName = getString("db_subnet_group_name")

	// Handle String Slice for Security Groups
	if val, ok := res.Metadata["vpc_security_group_ids"]; ok {
		if sl, ok := val.([]interface{}); ok {
			for _, s := range sl {
				if str, ok := s.(string); ok {
					rds.VpcSecurityGroupIds = append(rds.VpcSecurityGroupIds, str)
				}
			}
		} else if sl, ok := val.([]string); ok {
			rds.VpcSecurityGroupIds = sl
		}
	}

	rds.SkipFinalSnapshot = getBool("skip_final_snapshot")
	rds.PubliclyAccessible = getBool("publicly_accessible")
	rds.MultiAZ = getBool("multi_az")
	rds.BackupRetentionPeriod = getInt("backup_retention_period")

	return rds, nil
}

// MapRDSInstance maps an RDSInstance to a TerraformBlock
func MapRDSInstance(rds *database.RDSInstance) (*tfmapper.TerraformBlock, error) {
	if rds == nil {
		return nil, fmt.Errorf("rds instance is nil")
	}

	attributes := make(map[string]mapper.TerraformValue)

	attributes["identifier"] = strVal(rds.Name)
	attributes["engine"] = strVal(rds.Engine)
	attributes["engine_version"] = strVal(rds.EngineVersion)
	attributes["instance_class"] = strVal(rds.InstanceClass)
	attributes["allocated_storage"] = intVal(rds.AllocatedStorage)

	if rds.StorageType != "" {
		attributes["storage_type"] = strVal(rds.StorageType)
	}

	if rds.Username != "" {
		attributes["username"] = strVal(rds.Username)
	}

	if rds.Password != "" {
		attributes["password"] = strVal(rds.Password)
	}

	// Boolean values
	attributes["skip_final_snapshot"] = boolVal(rds.SkipFinalSnapshot)
	attributes["publicly_accessible"] = boolVal(rds.PubliclyAccessible)
	attributes["multi_az"] = boolVal(rds.MultiAZ)

	if rds.DBName != "" {
		attributes["db_name"] = strVal(rds.DBName)
	}

	if rds.SubnetGroupName != "" {
		attributes["db_subnet_group_name"] = strVal(rds.SubnetGroupName)
	}

	if len(rds.VpcSecurityGroupIds) > 0 {
		attributes["vpc_security_group_ids"] = listStrVal(rds.VpcSecurityGroupIds)
	}

	if rds.BackupRetentionPeriod > 0 {
		attributes["backup_retention_period"] = intVal(rds.BackupRetentionPeriod)
	}

	if rds.Tags != nil && len(rds.Tags) > 0 {
		tagsMap := make(map[string]mapper.TerraformValue)
		for _, tag := range rds.Tags {
			tagsMap[tag.Key] = strVal(tag.Value)
		}
		attributes["tags"] = mapper.TerraformValue{Map: tagsMap}
	}

	return &mapper.TerraformBlock{
		Kind:       "resource",
		Labels:     []string{"aws_db_instance", rds.Name},
		Attributes: attributes,
	}, nil
}

// Helper functions for creating TerraformValue

func strVal(s string) mapper.TerraformValue {
	return mapper.TerraformValue{String: &s}
}

func intVal(i int) mapper.TerraformValue {
	f := float64(i)
	return mapper.TerraformValue{Number: &f}
}

func boolVal(b bool) mapper.TerraformValue {
	return mapper.TerraformValue{Bool: &b}
}

func listStrVal(list []string) mapper.TerraformValue {
	vals := make([]mapper.TerraformValue, len(list))
	for i, s := range list {
		str := s // Create a new variable to take address of
		vals[i] = mapper.TerraformValue{String: &str}
	}
	return mapper.TerraformValue{List: vals}
}
