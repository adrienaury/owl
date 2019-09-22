package ldap

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Session in between commands
type Session struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
}

// Restore session from file
func Restore(file string) *Session {
	session := &Session{}

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.UnmarshalStrict(dat, session)
	if err != nil {
		log.Fatal(err)
	}

	return session
}

// Dump session into file
func (s *Session) Dump(file string) {
	dat, err := yaml.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(file, dat, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
