package models

// DependencyType represents types of resource dependencies (uses, depends_on, connects_to, references)
type DependencyType struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:text;uniqueIndex;not null" json:"name"`

	// Relationships
	ResourceDependencies []ResourceDependency `gorm:"foreignKey:DependencyTypeID" json:"resource_dependencies,omitempty"`
}

// TableName specifies the table name for GORM
func (DependencyType) TableName() string {
	return "dependency_types"
}
