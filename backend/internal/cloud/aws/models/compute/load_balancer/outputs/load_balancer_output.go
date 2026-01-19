package outputs

import (
	"time"
)

// LoadBalancerOutput represents AWS Load Balancer output/response data after creation
type LoadBalancerOutput struct {
	// AWS-generated identifiers
	ARN     string `json:"arn"` // e.g., "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-alb/..."
	ID      string `json:"id"`  // Same as ARN
	Name    string `json:"name"`
	DNSName string `json:"dns_name"` // Auto-generated DNS name
	ZoneID  string `json:"zone_id"`  // Route53 hosted zone ID

	// Configuration (from input)
	Type             string   `json:"type"`              // application or network
	Internal         bool     `json:"internal"`
	SecurityGroupIDs []string `json:"security_group_ids"`
	SubnetIDs        []string `json:"subnet_ids"`

	// AWS-specific output fields
	State       string    `json:"state"`        // active, provisioning, active_impaired, failed
	CreatedTime time.Time `json:"created_time"`
}
