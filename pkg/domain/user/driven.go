package user

// Backend ...
type Backend interface {
	ListUsers() (List, error)
	CreateUser(User) error
	DeleteUser(id string) error
}
