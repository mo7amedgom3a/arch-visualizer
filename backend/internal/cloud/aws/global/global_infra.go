package global

var REGIONS = []string{
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
	"eu-west-1",
	"eu-west-2",
	"eu-central-1",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-southeast-1",
}
// get all regions
func GetAllRegions() []string {
	return REGIONS
}
// get region by name
func GetRegionByName(name string) string {
	for _, region := range REGIONS {
		if region == name {
			return region
		}
	}
	return ""
}

// Availabiltiy Zones 
var AVAILABILITY_ZONES = map[string][]string{
	"us-east-1": {"us-east-1a", "us-east-1b", "us-east-1c"},
	"us-east-2": {"us-east-2a", "us-east-2b", "us-east-2c"},
	"us-west-1": {"us-west-1a", "us-west-1b", "us-west-1c"},
	"us-west-2": {"us-west-2a", "us-west-2b", "us-west-2c"},
	"eu-west-1": {"eu-west-1a", "eu-west-1b", "eu-west-1c"},
	"eu-west-2": {"eu-west-2a", "eu-west-2b", "eu-west-2c"},
	"eu-central-1": {"eu-central-1a", "eu-central-1b", "eu-central-1c"},
	"ap-northeast-1": {"ap-northeast-1a", "ap-northeast-1b", "ap-northeast-1c"},
	"ap-northeast-2": {"ap-northeast-2a", "ap-northeast-2b", "ap-northeast-2c"},
	"ap-southeast-1": {"ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"},
}

// get availability zones by region
func GetAvailabilityZonesByRegion(region string) []string {
	return AVAILABILITY_ZONES[region]
}
