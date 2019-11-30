package realm

// Storage ...
type Storage interface {
	CreateOrUpdateRealm(r Realm) error
	GetRealm(id string) (Realm, error)
	ListRealms() (List, error)
	DeleteRealm(id string) error
}
