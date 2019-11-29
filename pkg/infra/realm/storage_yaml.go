package realm

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrienaury/owl/pkg/domain/realm"
	"gopkg.in/yaml.v3"
)

// Version of the YAML strcuture
const Version string = "v1.beta"

// YAMLStructure of the file
type YAMLStructure struct {
	Version string      `yaml:"version"`
	Realms  []YAMLRealm `yaml:"realms,omitempty"`
}

// YAMLRealm defines how to store a realm in YAML format
type YAMLRealm struct {
	ID       string `yaml:"id"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
}

// NewYAMLStorage create a new YAML storage
func NewYAMLStorage() YAMLStorage {
	return YAMLStorage{}
}

// YAMLStorage provides storage in a local YAML file
type YAMLStorage struct{}

// CreateOrUpdateRealm in the local YAML file
func (s YAMLStorage) CreateOrUpdateRealm(r realm.Realm) error {
	structure, err := readFile()
	if err != nil {
		return err
	}

	list := structure.Realms

	result := list[:0]
	updated := false
	for _, ym := range list {
		if ym.ID == r.ID() {
			updated = true
			ym.URL = r.URL()
			ym.Username = r.Username()
		}
		result = append(result, ym)
	}

	if !updated {
		result = append(result, YAMLRealm{r.ID(), r.URL(), r.Username()})
	}

	return writeFile(&YAMLStructure{Version, result})
}

// GetRealm from the local YAML file
func (s YAMLStorage) GetRealm(id string) (realm.Realm, error) {
	structure, err := readFile()
	if err != nil {
		return nil, err
	}

	list := structure.Realms

	for _, ym := range list {
		if ym.ID == id {
			return realm.NewRealm(ym.ID, ym.URL, ym.Username), nil
		}
	}

	return nil, nil
}

// ListRealms contained in the local YAML file
func (s YAMLStorage) ListRealms() ([]realm.Realm, error) {
	structure, err := readFile()
	if err != nil {
		return nil, err
	}

	list := structure.Realms
	result := make([]realm.Realm, len(list))

	for _, ym := range list {
		m := realm.NewRealm(ym.ID, ym.URL, ym.Username)
		result = append(result, m)
	}

	return result, nil
}

func readFile() (*YAMLStructure, error) {
	structure := &YAMLStructure{
		Version: Version,
	}

	if _, err := os.Stat("realms.yaml"); os.IsNotExist(err) {
		return structure, nil
	}

	dat, err := ioutil.ReadFile("realms.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dat, structure)
	if err != nil {
		return nil, err
	}

	if structure.Version != Version {
		return nil, fmt.Errorf("invalid version in ./realms.yaml (%v)", structure.Version)
	}

	return structure, nil
}

func writeFile(structure *YAMLStructure) error {
	out, err := yaml.Marshal(structure)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("realms.yaml", out, 0640)
	if err != nil {
		return err
	}

	return nil
}
