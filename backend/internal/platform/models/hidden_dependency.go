package models

// HiddenDependency represents hidden/implicit dependencies between resources
type HiddenDependency struct {
	ID                  uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Provider            string `gorm:"type:varchar(20);not null;index" json:"provider"`
	ParentResourceType  string `gorm:"type:varchar(100);not null;index" json:"parent_resource_type"`
	ChildResourceType   string `gorm:"type:varchar(100);not null;index" json:"child_resource_type"`
	QuantityExpression  string `gorm:"type:varchar(255);default:'1'" json:"quantity_expression"`
	ConditionExpression string `gorm:"type:varchar(255)" json:"condition_expression,omitempty"`
	IsAttached          bool   `gorm:"default:true" json:"is_attached"`
	Description         string `gorm:"type:text" json:"description,omitempty"`
}

// TableName specifies the table name for GORM
func (HiddenDependency) TableName() string {
	return "hidden_dependencies"
}
