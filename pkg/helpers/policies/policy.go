package policies

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// Policies todo
type Policies struct {
	Version  string            `yaml:"version"`
	Policies map[string]Policy `yaml:"policies,omitempty"`
}

// A Policy configure how data is queried in LDAP
type Policy struct {
	Attributes map[string]Attribute `yaml:"attributes,omitempty"`
	Filter     Filter               `yaml:"filter"`
}

// Attribute of policy
type Attribute struct {
	Rules []Rule `yaml:"rules,omitempty"`
}

// Filter of policy
type Filter string

// Rule to apply to an attribute
type Rule string

var rules = map[Rule]func(string) string{
	"Capitalize": func(s string) string { return strings.Title(s) },
	"HashSha256": func(s string) string { return s },
}

// Get policies from file
func Get(file string) (*Policies, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("invalid policy file: %s", err)
	}

	p := &Policies{}
	err = yaml.UnmarshalStrict(dat, p)
	if err != nil {
		return nil, fmt.Errorf("invalid policy file: %s", err)
	}

	if p.Version != "v1.beta1" {
		return p, fmt.Errorf("invalid policy file version: %s", p.Version)
	}

	return p, nil
}

// GetAttributes todo
func (p *Policy) GetAttributes() []string {
	keys := make([]string, len(p.Attributes))
	i := 0
	for k := range p.Attributes {
		keys[i] = k
		i++
	}
	return keys
}

// GetFilter todo
func (p *Policy) GetFilter() string {
	return string(p.Filter)
}
