package models

import (
	"time"

	"github.com/google/uuid"
)

// Review represents a marketplace review
type Review struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TemplateID      uuid.UUID `gorm:"type:uuid;not null;index" json:"template_id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Rating          int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Title           string    `gorm:"type:varchar(255);not null" json:"title"`
	Content         string    `gorm:"type:text;not null" json:"content"`
	UseCase         *string   `gorm:"type:varchar(255)" json:"use_case,omitempty"`
	TeamSize        *string   `gorm:"type:varchar(50)" json:"team_size,omitempty"`
	DeploymentTime  *string   `gorm:"type:varchar(50)" json:"deployment_time,omitempty"`
	HelpfulCount    int       `gorm:"default:0" json:"helpful_count"`
	CreatorResponse *string   `gorm:"type:text" json:"creator_response,omitempty"`
	CreatedAt       time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:current_timestamp" json:"updated_at"`

	// Relationships
	Template Template `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
	User     User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName specifies the table name for GORM
func (Review) TableName() string {
	return "reviews"
}
