package diagram

import (
	"testing"
)

func TestConvertDBNameToIRType(t *testing.T) {
	tests := []struct {
		dbName   string
		expected string
	}{
		{"RouteTable", "route-table"},
		{"SecurityGroup", "security-group"},
		{"VPC", "vpc"},
		{"EC2", "ec2"},
		{"AutoScalingGroup", "auto-scaling-group"},
		{"InternetGateway", "internet-gateway"},
		{"NATGateway", "nat-gateway"},
		{"ElasticIP", "elastic-ip"},
		{"LoadBalancer", "load-balancer"},
		{"Lambda", "lambda"},
		{"S3", "s3"},
		{"EBS", "ebs"},
		{"RDS", "rds"},
		{"DynamoDB", "dynamodb"},
		{"Subnet", "subnet"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.dbName, func(t *testing.T) {
			result := convertDBNameToIRType(tt.dbName)
			if result != tt.expected {
				t.Errorf("convertDBNameToIRType(%q) = %q, want %q", tt.dbName, result, tt.expected)
			}
		})
	}
}
