package models

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a marketplace category for templates
type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Slug      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"slug"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`

	// Relationships
	Templates []Template `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"templates,omitempty"`
}

// TableName specifies the table name for GORM
func (Category) TableName() string {
	return "categories"
}
