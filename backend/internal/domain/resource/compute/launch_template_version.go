package compute

import "time"

// LaunchTemplateVersion represents a specific version of a launch template
type LaunchTemplateVersion struct {
	TemplateID    string
	VersionNumber int
	IsDefault     bool
	CreateTime    time.Time
	CreatedBy     *string // IAM user/role ARN
	TemplateData  *LaunchTemplate // Template configuration for this version
}
