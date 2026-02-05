package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ResourceUIState represents the UI-specific state for a resource node
type ResourceUIState struct {
	ID         uint           `gorm:"primary_key" json:"id"`
	ResourceID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"resource_id"`
	X          float64        `gorm:"not null;default:0" json:"x"`
	Y          float64        `gorm:"not null;default:0" json:"y"`
	Width      *float64       `json:"width"`
	Height     *float64       `json:"height"`
	Style      datatypes.JSON `gorm:"type:jsonb" json:"style"`
	Measured   datatypes.JSON `gorm:"type:jsonb" json:"measured"`
	Selected   bool           `gorm:"default:false" json:"selected"`
	Dragging   bool           `gorm:"default:false" json:"dragging"`
	Resizing   bool           `gorm:"default:false" json:"resizing"`
	Focusable  bool           `gorm:"default:true" json:"focusable"`
	Selectable bool           `gorm:"default:true" json:"selectable"`
	ZIndex     int            `gorm:"default:0" json:"z_index"`
	CreatedAt  time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (ResourceUIState) TableName() string {
	return "resource_ui_states"
}

// ProjectUIState represents the global UI state for a project view
type ProjectUIState struct {
	ID              uint           `gorm:"primary_key" json:"id"`
	ProjectID       uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"project_id"`
	Zoom            float64        `gorm:"default:1.0" json:"zoom"`
	ViewportX       float64        `gorm:"default:0" json:"viewport_x"`
	ViewportY       float64        `gorm:"default:0" json:"viewport_y"`
	SelectedNodeIDs datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"selected_node_ids"`
	SelectedEdgeIDs datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"selected_edge_ids"`
	CreatedAt       time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"default:now()" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (ProjectUIState) TableName() string {
	return "project_ui_states"
}
