package realm

// Realm ...
type Realm interface {
	ID() string
	URL() string
	Username() string
}

type realm struct {
	id       string
	url      string
	username string
}

// NewRealm ...
func NewRealm(id, url, username string) Realm {
	return realm{
		id:       id,
		url:      url,
		username: username,
	}
}

func (r realm) ID() string       { return r.id }
func (r realm) URL() string      { return r.url }
func (r realm) Username() string { return r.username }
