package models

// PricingComponent represents pricing breakdown per component (e.g., per-hour, per-GB, per-request)
type PricingComponent struct {
	ID               uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourcePricingID uint   `gorm:"not null;index" json:"resource_pricing_id"`
	ComponentName    string  `gorm:"type:text;not null" json:"component_name"`
	Model            string  `gorm:"type:text;not null;check:model IN ('per_hour','per_gb','per_request','one_time','tiered','percentage')" json:"model"`
	Unit             string  `gorm:"type:text;not null" json:"unit"`
	Quantity         float64 `gorm:"type:numeric(14,4);not null" json:"quantity"`
	UnitRate         float64 `gorm:"type:numeric(14,6);not null" json:"unit_rate"`
	Subtotal         float64 `gorm:"type:numeric(14,4);not null" json:"subtotal"`
	Currency         string  `gorm:"type:text;not null;check:currency IN ('USD','EUR','GBP')" json:"currency"`

	// Relationships
	ResourcePricing ResourcePricing `gorm:"foreignKey:ResourcePricingID;constraint:OnDelete:CASCADE" json:"resource_pricing,omitempty"`
}

// TableName specifies the table name for GORM
func (PricingComponent) TableName() string {
	return "pricing_components"
}
