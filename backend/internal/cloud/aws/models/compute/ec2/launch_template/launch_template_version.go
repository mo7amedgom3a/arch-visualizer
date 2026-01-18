package launch_template

import "time"

// LaunchTemplateVersion represents a specific version of a launch template
type LaunchTemplateVersion struct {
	TemplateID    string    `json:"template_id"`    // Launch template ID
	VersionNumber int       `json:"version_number"` // Version number
	IsDefault     bool      `json:"is_default"`      // Whether this is the default version
	CreateTime    time.Time `json:"create_time"`    // When this version was created
	CreatedBy     *string   `json:"created_by,omitempty"` // IAM user/role ARN
	TemplateData  *LaunchTemplate `json:"template_data,omitempty"` // Template configuration for this version
}
