package models

import "github.com/google/uuid"

// TemplateIACFormat represents the template_iac_formats join table
type TemplateIACFormat struct {
	TemplateID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"template_id"`
	IACFormatID uuid.UUID `gorm:"type:uuid;primaryKey" json:"iac_format_id"`

	// Relationships
	Template  Template  `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
	IACFormat IACFormat `gorm:"foreignKey:IACFormatID;constraint:OnDelete:CASCADE" json:"iac_format,omitempty"`
}

// TableName specifies the table name for GORM
func (TemplateIACFormat) TableName() string {
	return "template_iac_formats"
}
