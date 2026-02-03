package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ProjectVariable represents a Terraform/IaC input variable for the project
type ProjectVariable struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProjectID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"projectId"`
	Name         string         `gorm:"not null" json:"name"`
	Type         string         `gorm:"not null" json:"type"` // string, number, bool, list(string), map(string)
	Description  string         `json:"description"`
	DefaultValue datatypes.JSON `json:"defaultValue"` // Store as JSONB
	Sensitive    bool           `gorm:"default:false" json:"sensitive"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName overrides the table name used by User to `project_variables`
func (ProjectVariable) TableName() string {
	return "project_variables"
}

// UnmarshalDefaultValue unmarshals the DefaultValue JSON into a Go interface{}
func (pv *ProjectVariable) UnmarshalDefaultValue() (interface{}, error) {
	if len(pv.DefaultValue) == 0 {
		return nil, nil
	}
	var val interface{}
	err := json.Unmarshal(pv.DefaultValue, &val)
	return val, err
}
