package models

// ResourceType represents a cloud-specific resource type (EC2, Lambda, S3, RDS, VPC, etc.)
type ResourceType struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"type:text;not null" json:"name"`
	CloudProvider string `gorm:"type:text;not null" json:"cloud_provider"`
	CategoryID    *uint  `gorm:"index" json:"category_id,omitempty"`
	KindID        *uint  `gorm:"index" json:"kind_id,omitempty"`
	IsRegional    bool  `gorm:"default:true" json:"is_regional"`
	IsGlobal      bool  `gorm:"default:false" json:"is_global"`

	// Relationships
	Category          *ResourceCategory      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Kind              *ResourceKind          `gorm:"foreignKey:KindID" json:"kind,omitempty"`
	Resources         []Resource             `gorm:"foreignKey:ResourceTypeID" json:"resources,omitempty"`
	Constraints       []ResourceConstraint   `gorm:"foreignKey:ResourceTypeID" json:"constraints,omitempty"`
	ServiceTypePricing []ServiceTypePricing   `gorm:"foreignKey:ResourceTypeID" json:"service_type_pricing,omitempty"`
}

// TableName specifies the table name for GORM
func (ResourceType) TableName() string {
	return "resource_types"
}
