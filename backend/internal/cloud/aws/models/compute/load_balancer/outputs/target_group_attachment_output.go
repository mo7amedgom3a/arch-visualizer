package outputs

// TargetGroupAttachmentOutput represents AWS Target Group Attachment output/response data
type TargetGroupAttachmentOutput struct {
	TargetGroupARN   string  `json:"target_group_arn"`
	TargetID         string  `json:"target_id"`
	Port             *int    `json:"port,omitempty"`
	AvailabilityZone *string `json:"availability_zone,omitempty"`
	HealthStatus     string  `json:"health_status"` // healthy, unhealthy, initial, draining
	State            string  `json:"state"`        // initial, healthy, unhealthy, unused, draining, unavailable
}
