package resource
type CloudProvider string
const (
	AWS CloudProvider = "aws"
	Azure CloudProvider = "azure"
	GCP CloudProvider = "gcp"
)
type AvailabilityZone struct {
    Name   string // us-east-1a
    Region string // us-east-1
}
type Region struct {
	Id string 
	Provider CloudProvider
	Name string
	AZs []CloudProvider
}