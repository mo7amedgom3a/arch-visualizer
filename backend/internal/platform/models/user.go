package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string         `gorm:"type:text;uniqueIndex;not null" json:"email"`
	Name      *string        `gorm:"type:text" json:"name,omitempty"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Projects []Project `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"projects,omitempty"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
