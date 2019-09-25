package policies

import (
	"fmt"
	"io/ioutil"

	"github.com/adrienaury/owl/configs"
	"github.com/adrienaury/owl/pkg/helpers/defaults"
	"gopkg.in/yaml.v3"
)

// Version todo
const Version string = "v1.beta1"

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

	if p.Version != Version {
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

func init() {
	var err error
	if !defaults.Exists(defaults.Policies) {
		err = ioutil.WriteFile(defaults.Policies, []byte(configs.Policies), 0644)
	} else {
		_, err = Get(defaults.Policies)
	}

	if err != nil {
		// TODO
		fmt.Println(err)
	}
}
