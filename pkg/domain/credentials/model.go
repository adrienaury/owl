package credentials

// Credentials object contains information on how to connect to a realm.
type Credentials interface {
	URL() string
	Username() string
	Password() string
}

type credentials struct {
	url      string
	username string
	password string
}

// NewCredentials create a new credentials object.
func NewCredentials(url, username, password string) Credentials {
	return credentials{
		url:      url,
		username: username,
		password: password,
	}
}

func (r credentials) URL() string      { return r.url }
func (r credentials) Username() string { return r.username }
func (r credentials) Password() string { return r.password }
