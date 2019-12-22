package realm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

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
	Policy   string `yaml:"policy,omitempty"`
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
	return path.Join(s.localDir, "realms.yaml")
}

// CreateOrUpdateRealm in the local YAML file
func (s YAMLStorage) CreateOrUpdateRealm(r realm.Realm) error {
	structure, err := s.readFile()
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
			ym.Policy = r.Policy()
		}
		result = append(result, ym)
	}

	if !updated {
		result = append(result, YAMLRealm{r.ID(), r.URL(), r.Username(), r.Policy()})
	}

	return s.writeFile(&YAMLStructure{Version, result})
}

// GetRealm from the local YAML file
func (s YAMLStorage) GetRealm(id string) (realm.Realm, error) {
	structure, err := s.readFile()
	if err != nil {
		return nil, err
	}

	list := structure.Realms

	for _, ym := range list {
		if ym.ID == id {
			return realm.NewRealmWithPolicy(ym.ID, ym.URL, ym.Username, ym.Policy), nil
		}
	}

	return nil, nil
}

// DeleteRealm remove the realm with id from the local storage.
func (s YAMLStorage) DeleteRealm(id string) error {
	structure, err := s.readFile()
	if err != nil {
		return err
	}

	list := structure.Realms

	result := list[:0]
	for _, ym := range list {
		if ym.ID != id {
			result = append(result, ym)
		}
	}

	return s.writeFile(&YAMLStructure{Version, result})
}

// ListRealms contained in the local YAML file
func (s YAMLStorage) ListRealms() (realm.List, error) {
	structure, err := s.readFile()
	if err != nil {
		return nil, err
	}

	list := structure.Realms
	result := make([]realm.Realm, 0, len(list))

	for _, ym := range list {
		m := realm.NewRealmWithPolicy(ym.ID, ym.URL, ym.Username, ym.Policy)
		result = append(result, m)
	}

	return realm.NewList(result), nil
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

func (s YAMLStorage) writeFile(structure *YAMLStructure) error {
	out, err := yaml.Marshal(structure)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.getFileName(), out, 0640)
	if err != nil {
		return err
	}

	return nil
}
