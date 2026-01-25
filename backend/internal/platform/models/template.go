package models

import (
	"time"

	"github.com/google/uuid"
)

// Template represents a marketplace template
type Template struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title            string    `gorm:"type:varchar(255);not null" json:"title"`
	Description      string    `gorm:"type:text;not null" json:"description"`
	CategoryID       uuid.UUID `gorm:"type:uuid;not null;index" json:"category_id"`
	CloudProvider    string    `gorm:"type:varchar(50);not null;check:cloud_provider IN ('AWS','Azure','GCP','Multi-Cloud')" json:"cloud_provider"`
	Rating           float64   `gorm:"type:decimal(3,2);default:0;check:rating >= 0 AND rating <= 5" json:"rating"`
	ReviewCount      int       `gorm:"default:0" json:"review_count"`
	Downloads        int       `gorm:"default:0" json:"downloads"`
	Price            float64   `gorm:"type:decimal(10,2);default:0" json:"price"`
	IsSubscription   bool      `gorm:"default:false" json:"is_subscription"`
	SubscriptionPrice *float64 `gorm:"type:decimal(10,2)" json:"subscription_price,omitempty"`
	EstimatedCostMin float64   `gorm:"type:decimal(10,2);not null" json:"estimated_cost_min"`
	EstimatedCostMax float64   `gorm:"type:decimal(10,2);not null" json:"estimated_cost_max"`
	AuthorID         uuid.UUID `gorm:"type:uuid;not null;index" json:"author_id"`
	ImageURL         *string   `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	IsPopular        bool      `gorm:"default:false" json:"is_popular"`
	IsNew            bool      `gorm:"default:false" json:"is_new"`
	LastUpdated      time.Time `gorm:"default:current_timestamp" json:"last_updated"`
	Resources        int       `gorm:"default:0" json:"resources"`
	DeploymentTime   *string   `gorm:"type:varchar(50)" json:"deployment_time,omitempty"`
	Regions          *string   `gorm:"type:text" json:"regions,omitempty"`
	CreatedAt        time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:current_timestamp" json:"updated_at"`

	// Relationships
	Category            Category             `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"category,omitempty"`
	Author              User                 `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"author,omitempty"`
	Technologies        []Technology         `gorm:"many2many:template_technologies" json:"technologies,omitempty"`
	IACFormats          []IACFormat          `gorm:"many2many:template_iac_formats" json:"iac_formats,omitempty"`
	ComplianceStandards []ComplianceStandard `gorm:"many2many:template_compliance;joinForeignKey:TemplateID;joinReferences:ComplianceID" json:"compliance_standards,omitempty"`
	UseCases            []TemplateUseCase    `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"use_cases,omitempty"`
	Features            []TemplateFeature    `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"features,omitempty"`
	Components          []TemplateComponent  `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"components,omitempty"`
	Reviews             []Review             `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
}

// TableName specifies the table name for GORM
func (Template) TableName() string {
	return "templates"
}
