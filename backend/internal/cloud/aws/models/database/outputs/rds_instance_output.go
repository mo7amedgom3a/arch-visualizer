package outputs

// RDSInstanceOutput represents the output attributes of an AWS RDS Instance
type RDSInstanceOutput struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Endpoint string `json:"endpoint"`
	ARN      string `json:"arn"`
}
