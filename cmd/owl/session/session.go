package session

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// Session store context in between commands
type Session struct {
	Version string `yaml:"version"`
	Realm   string `yaml:"realm"`
	Unit    string `yaml:"unit"`

	sourcefile string
}

// NewSession creates a new session
func NewSession(sourcefile string) *Session {
	return &Session{Version: "v1.beta1", sourcefile: sourcefile}
}

// Restore session from file
func (s *Session) Restore() (*Session, error) {
	dat, err := ioutil.ReadFile(s.sourcefile)
	if os.IsNotExist(err) {
		s.Dump()
	} else if err != nil {
		return s, fmt.Errorf("invalid session file: %s", err)
	}

	err = yaml.Unmarshal(dat, s)
	if err != nil {
		return s, fmt.Errorf("invalid session file: %s", err)
	}

	if s.Version != "v1.beta1" {
		return s, fmt.Errorf("invalid session file version: %s", s.Version)
	}

	return s, nil
}

// Dump session into file
func (s *Session) Dump() error {
	if s == nil {
		return fmt.Errorf("no session file defined")
	}

	dat, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.sourcefile, dat, 0644)
	if err != nil {
		return err
	}

	return nil
}
