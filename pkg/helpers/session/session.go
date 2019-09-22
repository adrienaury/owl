package session

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Session store context in between commands
type Session struct {
	Version  string `yaml:"version"`
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
}

// NewSession creates a new session
func NewSession() *Session {
	return &Session{Version: "v1.beta1"}
}

// Restore session from file
func (s *Session) Restore(file string) (*Session, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return s, fmt.Errorf("invalid session file: %s", err)
	}

	err = yaml.UnmarshalStrict(dat, s)
	if err != nil {
		return s, fmt.Errorf("invalid session file: %s", err)
	}

	if s.Version != "v1.beta1" {
		return s, fmt.Errorf("invalid session file version: %s", s.Version)
	}

	return s, nil
}

// Dump session into file
func (s *Session) Dump(file string) error {
	dat, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, dat, 0644)
	if err != nil {
		return err
	}

	return nil
}
