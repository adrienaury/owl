package user

// Backend interface.
type Backend interface {
	ListUsers() (List, error)
	GetUser(id string) (User, error)
	CreateUser(User) error
	UpdateUser(User) error
	AppendUser(User) error
	RemoveUser(User) error
	DeleteUser(id string) error
}
