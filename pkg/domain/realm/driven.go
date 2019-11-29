package realm

// Storage ...
type Storage interface {
	CreateOrUpdateRealm(r Realm) error
	GetRealm(id string) (Realm, error)
	ListRealms() ([]Realm, error)
}
