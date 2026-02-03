package models

import (
	"time"

	"github.com/google/uuid"
)

// ProjectOutput represents a Terraform/IaC output definition for the project
type ProjectOutput struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProjectID   uuid.UUID `gorm:"type:uuid;not null;index" json:"projectId"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Value       string    `gorm:"not null" json:"value"` // The expression, e.g. "${aws_instance.main.public_ip}"
	Sensitive   bool      `gorm:"default:false" json:"sensitive"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName overrides the table name used by User to `project_outputs`
func (ProjectOutput) TableName() string {
	return "project_outputs"
}
