package credentials

// Storage ...
type Storage interface {
	SetCredentials(c Credentials) error
	GetCredentials(url string, user string) (Credentials, error)
}

// Backend ...
type Backend interface {
	TestCredentials(c Credentials) (bool, error)
}
