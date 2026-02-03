package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Email      string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Auth0ID    string    `gorm:"type:varchar(255);unique;not null;index" json:"auth0_id"`
	Avatar     *string   `gorm:"type:varchar(500)" json:"avatar,omitempty"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`

	// Relationships
	Projects  []Project  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"projects,omitempty"`
	Templates []Template `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"templates,omitempty"`
	Reviews   []Review   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
