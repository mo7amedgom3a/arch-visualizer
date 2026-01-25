package models

import (
	"time"

	"github.com/google/uuid"
)

// Technology represents a marketplace technology tag
type Technology struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Slug      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"slug"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`

	// Relationships
	Templates []Template `gorm:"many2many:template_technologies" json:"templates,omitempty"`
}

// TableName specifies the table name for GORM
func (Technology) TableName() string {
	return "technologies"
}
