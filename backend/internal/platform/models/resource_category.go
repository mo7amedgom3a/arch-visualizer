package models

// ResourceCategory represents a resource category (Compute, Networking, Storage, etc.)
type ResourceCategory struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:text;uniqueIndex;not null" json:"name"`

	// Relationships
	ResourceTypes []ResourceType `gorm:"foreignKey:CategoryID" json:"resource_types,omitempty"`
	ServicePricing []ServicePricing `gorm:"foreignKey:CategoryID" json:"service_pricing,omitempty"`
}

// TableName specifies the table name for GORM
func (ResourceCategory) TableName() string {
	return "resource_categories"
}
