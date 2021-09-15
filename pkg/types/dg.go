package types

type ResourceGroupTree struct {
	ChildGroups []*ResourceGroupTree
	Options     []ResourceGroupOption
	DontCreate  bool
}
