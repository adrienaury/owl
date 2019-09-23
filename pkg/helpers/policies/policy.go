package policies

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Policies todo
type Policies struct {
	Version  string            `yaml:"version"`
	Policies map[string]Policy `yaml:"policies,omitempty"`
}

// A Policy configure how data is queried in LDAP
type Policy struct {
	Attributes []Attribute `yaml:"attributes"`
	Filter     Filter      `yaml:"filter"`
}

// Attribute of policy
type Attribute struct {
	Name  string   `yaml:"name"`
	Rules []string `yaml:"rules,omitempty"`
}

// Filter of policy
type Filter string

// Get policies from file
func Get(file string) (*Policies, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("invalid policy file: %s", err)
	}

	p := &Policies{}
	err = yaml.Unmarshal(dat, p)
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
	for i, k := range p.Attributes {
		keys[i] = k.Name
	}
	return keys
}

// GetFilter todo
func (p *Policy) GetFilter() string {
	return string(p.Filter)
}
