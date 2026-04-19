package database

import "errors"

// RDSInstance represents a cloud-agnostic database instance
type RDSInstance struct {
	ID                    string
	Name                  string
	Engine                string
	EngineVersion         string
	InstanceClass         string
	AllocatedStorage      int
	StorageType           string
	Username              string
	Password              string
	DBName                string
	SubnetGroupName       string
	VpcSecurityGroupIds   []string
	SkipFinalSnapshot     bool
	PubliclyAccessible    bool
	MultiAZ               bool
	BackupRetentionPeriod int
	Tags                  map[string]string

	// Output fields
	Endpoint string
	Port     int
	ARN      string
}

// Validate performs domain-level validation
func (i *RDSInstance) Validate() error {
	if i.Name == "" {
		return errors.New("rds instance name is required")
	}
	if i.Engine == "" {
		return errors.New("engine is required")
	}
	if i.InstanceClass == "" {
		return errors.New("instance class is required")
	}
	return nil
}
