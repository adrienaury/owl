package policy

// Policy about an object type.
type Policy interface {
	Name() string
	Objects() map[string]Object
}

type policy struct {
	name    string
	objects map[string]Object
}

// NewPolicy create a new policy.
func NewPolicy(name string, objects []Object) Policy {
	objmap := map[string]Object{}
	for _, obj := range objects {
		objmap[obj.Name()] = obj
	}
	return policy{
		name:    name,
		objects: objmap,
	}
}

func (p policy) Name() string               { return p.name }
func (p policy) Objects() map[string]Object { return p.objects }

// Object policy.
type Object interface {
	Name() string
	BackendObject() string
	BackendFields() map[string]string
}

type object struct {
	name    string
	backend string
	fields  map[string]string
}

// NewObject create a new object policy.
func NewObject(name, backend string, fields map[string]string) Object {
	return object{
		name:    name,
		backend: backend,
		fields:  fields,
	}
}

func (p object) Name() string                     { return p.name }
func (p object) BackendObject() string            { return p.backend }
func (p object) BackendFields() map[string]string { return p.fields }

var defaultUnitPolicy = NewObject(
	"unit",
	"organizationalUnit",
	map[string]string{
		"id":          "ou",
		"description": "description",
	},
)

var defaultUserPolicy = NewObject(
	"user",
	"inetOrgPerson",
	map[string]string{
		"id":         "cn",
		"firstnames": "givenName",
		"lastnames":  "sn",
		"emails":     "mail",
		"password":   "userPassword",
	},
)

var defaultGroupPolicy = NewObject(
	"group",
	"groupOfUniqueNames",
	map[string]string{
		"id":      "cn",
		"members": "uniqueMember",
	},
)

var defaultPolicy = NewPolicy("default", []Object{
	defaultUnitPolicy,
	defaultUserPolicy,
	defaultGroupPolicy,
})
