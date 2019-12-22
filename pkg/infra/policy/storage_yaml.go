package policy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/adrienaury/owl/pkg/domain/policy"
	"gopkg.in/yaml.v3"
)

// Version of the YAML strcuture
const Version string = "v1.beta"

// YAMLStructure of the file
type YAMLStructure struct {
	Version  string       `yaml:"version"`
	Policies []YAMLPolicy `yaml:"policies,omitempty"`
}

// YAMLPolicy defines how to store a policy in YAML format
type YAMLPolicy struct {
	Name    string       `yaml:"name"`
	Objects []YAMLObject `yaml:"objects"`
}

// YAMLObject defines how to store an object policy in YAML format
type YAMLObject struct {
	Name    string      `yaml:"name"`
	Backend string      `yaml:"backend"`
	Fields  []YAMLField `yaml:"fields"`
}

// YAMLField defines how to store an field policy in YAML format
type YAMLField struct {
	Name        string `yaml:"name"`
	Backend     string `yaml:"backend"`
	Multivalued bool   `yaml:"multivalued"`
}

// NewYAMLStorage create a new YAML storage
func NewYAMLStorage(localDir string) YAMLStorage {
	return YAMLStorage{localDir}
}

// YAMLStorage provides storage in a local YAML file
type YAMLStorage struct {
	localDir string
}

func (s YAMLStorage) getFileName() string {
	return path.Join(s.localDir, "policies.yaml")
}

// GetPolicy from the local YAML file
func (s YAMLStorage) GetPolicy(name string) (policy.Policy, error) {
	structure, err := s.readFile()
	if err != nil {
		return nil, err
	}

	list := structure.Policies

	for _, ym := range list {
		if ym.Name == name {
			return policy.NewPolicy(ym.Name, s.objectFromYAML(ym)), nil
		}
	}

	return nil, nil
}

func (s YAMLStorage) objectFromYAML(ym YAMLPolicy) []policy.Object {
	result := []policy.Object{}
	for _, obj := range ym.Objects {
		result = append(result, policy.NewObject(obj.Name, obj.Backend, s.fieldsFromYAML(obj.Fields)))
	}
	return result
}

func (s YAMLStorage) fieldsFromYAML(fields []YAMLField) map[string]string {
	result := map[string]string{}
	for _, field := range fields {
		result[field.Name] = field.Backend
	}
	return result
}

func (s YAMLStorage) readFile() (*YAMLStructure, error) {
	structure := &YAMLStructure{
		Version: Version,
	}

	if _, err := os.Stat(s.getFileName()); os.IsNotExist(err) {
		return structure, nil
	}

	dat, err := ioutil.ReadFile(s.getFileName())
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dat, structure)
	if err != nil {
		return nil, err
	}

	if structure.Version != Version {
		return nil, fmt.Errorf("invalid version in %v (%v)", s.getFileName(), structure.Version)
	}

	return structure, nil
}

/*func (s YAMLStorage) writeFile(structure *YAMLStructure) error {
	out, err := yaml.Marshal(structure)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.getFileName(), out, 0640)
	if err != nil {
		return err
	}

	return nil
}*/
