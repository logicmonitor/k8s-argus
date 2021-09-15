package enums

//go:generate go run resource_generator.go
// BuilderAction action to send to build while modifying system.categories property
type BuilderAction uint32

const (
	// Add add action
	Add BuilderAction = iota
	// Delete delete action
	Delete
)

// Is equality check to add if conditions for callers
func (ca *BuilderAction) Is(action BuilderAction) bool {
	return *ca == action
}
