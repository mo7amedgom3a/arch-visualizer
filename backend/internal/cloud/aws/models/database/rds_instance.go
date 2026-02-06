package database

import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"

// RDSInstance represents an AWS RDS Database Instance
type RDSInstance struct {
	Name                  string        `json:"name"`
	Engine                string        `json:"engine"` // mysql, postgres, mariadb, oracle-ee, sqlserver-ex, etc.
	EngineVersion         string        `json:"engine_version"`
	InstanceClass         string        `json:"instance_class"`         // db.t3.micro, db.m5.large, etc.
	AllocatedStorage      int           `json:"allocated_storage"`      // In GiB
	ReplicateSourceDB     string        `json:"replicate_source_db,omitempty"`
	StorageType           string        `json:"storage_type,omitempty"` // gp2, gp3, io1, standard
	Username              string        `json:"username,omitempty"`
	Password              string        `json:"password,omitempty"`
	DBName                string        `json:"db_name,omitempty"`
	SubnetGroupName       string        `json:"subnet_group_name,omitempty"`
	VpcSecurityGroupIds   []string      `json:"vpc_security_group_ids,omitempty"`
	SkipFinalSnapshot     bool          `json:"skip_final_snapshot,omitempty"`
	PubliclyAccessible    bool          `json:"publicly_accessible,omitempty"`
	MultiAZ               bool          `json:"multi_az,omitempty"`
	BackupRetentionPeriod int           `json:"backup_retention_period,omitempty"` // In days
	Tags                  []configs.Tag `json:"tags,omitempty"`
}
