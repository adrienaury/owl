package policies

import (
	"fmt"
	"io/ioutil"

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
