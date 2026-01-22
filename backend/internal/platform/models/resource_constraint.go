package models

// ResourceConstraint represents validation rules and constraints for resource types
type ResourceConstraint struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceTypeID uint  `gorm:"not null;index" json:"resource_type_id"`
	ConstraintType string `gorm:"type:text;not null" json:"constraint_type"`
	ConstraintValue string `gorm:"type:text;not null" json:"constraint_value"`

	// Relationships
	ResourceType ResourceType `gorm:"foreignKey:ResourceTypeID" json:"resource_type,omitempty"`
}

// TableName specifies the table name for GORM
func (ResourceConstraint) TableName() string {
	return "resource_constraints"
}
