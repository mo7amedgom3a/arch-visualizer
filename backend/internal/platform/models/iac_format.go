package models

import (
	"time"

	"github.com/google/uuid"
)

// IACFormat represents supported IaC formats in marketplace
type IACFormat struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Slug      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"slug"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`

	// Relationships
	Templates []Template `gorm:"many2many:template_iac_formats" json:"templates,omitempty"`
}

// TableName specifies the table name for GORM
func (IACFormat) TableName() string {
	return "iac_formats"
}
