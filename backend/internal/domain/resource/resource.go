package resource

type ResourceType struct {
	ID string
	Name string
	Category string
	Kind string
	IsRegional bool
	IsGlobal bool
}
type Resource struct {
    ID           string
    Name         string
    Type         ResourceType
    Provider     CloudProvider
    Region       string
    ParentID     *string
    DependsOn    []string
}