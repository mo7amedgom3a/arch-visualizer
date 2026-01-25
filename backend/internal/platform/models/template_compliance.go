package models

import "github.com/google/uuid"

// TemplateCompliance represents the template_compliance join table
type TemplateCompliance struct {
	TemplateID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"template_id"`
	ComplianceID uuid.UUID `gorm:"type:uuid;primaryKey" json:"compliance_id"`

	// Relationships
	Template   Template           `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
	Compliance ComplianceStandard `gorm:"foreignKey:ComplianceID;constraint:OnDelete:CASCADE" json:"compliance,omitempty"`
}

// TableName specifies the table name for GORM
func (TemplateCompliance) TableName() string {
	return "template_compliance"
}
