package models

import (
	"time"

	"github.com/google/uuid"
)

// ComplianceStandard represents a compliance standard in marketplace
type ComplianceStandard struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Slug      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"slug"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`

	// Relationships
	Templates []Template `gorm:"many2many:template_compliance" json:"templates,omitempty"`
}

// TableName specifies the table name for GORM
func (ComplianceStandard) TableName() string {
	return "compliance_standards"
}
