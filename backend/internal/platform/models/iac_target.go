package models

// IACTarget represents an IaC tool (Terraform, Pulumi, CDK)
type IACTarget struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:text;uniqueIndex;not null" json:"name"`

	// Relationships
	Projects []Project `gorm:"foreignKey:InfraToolID" json:"projects,omitempty"`
}

// TableName specifies the table name for GORM
func (IACTarget) TableName() string {
	return "iac_targets"
}
