package models

// ResourceKind represents a resource kind (VirtualMachine, Container, Function, etc.)
type ResourceKind struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:text;uniqueIndex;not null" json:"name"`

	// Relationships
	ResourceTypes []ResourceType `gorm:"foreignKey:KindID" json:"resource_types,omitempty"`
}

// TableName specifies the table name for GORM
func (ResourceKind) TableName() string {
	return "resource_kinds"
}
