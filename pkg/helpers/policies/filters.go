package policies

// NamedFilters can be used as a filter in a request
var NamedFilters = map[string]string{
	"all":   "(objectClass=*)",
	"users": "(userPassword=*)",
}

// NamedFiltersAttributes configure what attributes to show for each filter
var NamedFiltersAttributes = map[string][]string{
	"all":   []string{"dn", "objectClass"},
	"users": []string{"dn", "cn", "description"},
}
