package compute

// ComputeResource represents any compute resource in the domain
// This interface allows for polymorphic handling of compute resources
type ComputeResource interface {
	GetID() string
	GetName() string
	GetSubnetID() string
	Validate() error
}

// Ensure Instance implements ComputeResource
var _ ComputeResource = (*Instance)(nil)

// Implement ComputeResource for Instance
func (i *Instance) GetID() string    { return i.ID }
func (i *Instance) GetName() string  { return i.Name }
func (i *Instance) GetSubnetID() string { return i.SubnetID }

// Ensure LaunchTemplate implements ComputeResource
var _ ComputeResource = (*LaunchTemplate)(nil)

// Implement ComputeResource for LaunchTemplate
func (lt *LaunchTemplate) GetID() string {
	return lt.ID
}

func (lt *LaunchTemplate) GetName() string {
	if lt.Name != "" {
		return lt.Name
	}
	if lt.NamePrefix != nil {
		return *lt.NamePrefix
	}
	return ""
}

func (lt *LaunchTemplate) GetSubnetID() string {
	// Launch Templates don't have a subnet ID directly
	// They are used by ASGs which handle subnet selection
	return ""
}
