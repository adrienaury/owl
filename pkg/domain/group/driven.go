package group

// Backend interface.
type Backend interface {
	ListGroups() (List, error)
	GetGroup(id string) (Group, error)
	CreateGroup(Group) error
	UpdateGroup(Group) error
	DeleteGroup(id string) error
	AppendGroup(Group) error
	RemoveGroup(Group) error
}
