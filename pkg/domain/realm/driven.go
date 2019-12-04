package realm

// Storage interface for storing realms.
type Storage interface {
	CreateOrUpdateRealm(r Realm) error
	GetRealm(id string) (Realm, error)
	ListRealms() (List, error)
	DeleteRealm(id string) error
}
