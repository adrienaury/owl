package credentials

// Storage interface for storing credentials.
type Storage interface {
	SetCredentials(c Credentials) error
	GetCredentials(url string, user string) (Credentials, error)
	RemoveCredentials(url string, username string) error
}

// Backend interface.
type Backend interface {
	TestCredentials(c Credentials) (bool, error)
	OpenConnection(c Credentials) error
}
