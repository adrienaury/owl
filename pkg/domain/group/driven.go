package group

// Backend ...
type Backend interface {
	ListGroups() (List, error)
	GetGroup(id string) (Group, error)
	CreateGroup(Group) error
	DeleteGroup(id string) error
	AddToGroup(ids ...string) error
}
