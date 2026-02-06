package database

import (
	awsconfigs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsdatabase "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/database"
	awsdatabaseoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/database/outputs"
	domaindatabase "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/database"
)

// FromDomainRDSInstance converts a domain RDSInstance to an AWS RDSInstance
func FromDomainRDSInstance(d *domaindatabase.RDSInstance) *awsdatabase.RDSInstance {
	if d == nil {
		return nil
	}

	tags := make([]awsconfigs.Tag, 0, len(d.Tags))
	for k, v := range d.Tags {
		tags = append(tags, awsconfigs.Tag{Key: k, Value: v})
	}

	return &awsdatabase.RDSInstance{
		Name:                  d.Name,
		Engine:                d.Engine,
		EngineVersion:         d.EngineVersion,
		InstanceClass:         d.InstanceClass,
		AllocatedStorage:      d.AllocatedStorage,
		StorageType:           d.StorageType,
		Username:              d.Username,
		Password:              d.Password,
		DBName:                d.DBName,
		SubnetGroupName:       d.SubnetGroupName,
		VpcSecurityGroupIds:   d.VpcSecurityGroupIds,
		SkipFinalSnapshot:     d.SkipFinalSnapshot,
		PubliclyAccessible:    d.PubliclyAccessible,
		MultiAZ:               d.MultiAZ,
		BackupRetentionPeriod: d.BackupRetentionPeriod,
		Tags:                  tags,
	}
}

// ToDomainRDSInstanceFromOutput converts an AWS RDSInstanceOutput to a domain RDSInstance
func ToDomainRDSInstanceFromOutput(o *awsdatabaseoutputs.RDSInstanceOutput) *domaindatabase.RDSInstance {
	if o == nil {
		return nil
	}

	return &domaindatabase.RDSInstance{
		ID:       o.ID,
		Endpoint: o.Address,
		Port:     o.Port,
		ARN:      o.ARN,
		// Typically output doesn't contain all config fields unless we fetch them from Describe,
		// but for now we populate what we have in output
	}
}
