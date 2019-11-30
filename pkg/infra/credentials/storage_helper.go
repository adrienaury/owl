package credentials

import (
	"runtime"

	"github.com/adrienaury/owl/pkg/domain/credentials"
	"github.com/docker/docker-credential-helpers/client"
	backend "github.com/docker/docker-credential-helpers/credentials"
)

// NewHelperStorage create a new YAML storage
func NewHelperStorage() HelperStorage {
	var nativeStore client.ProgramFunc
	switch runtime.GOOS {
	case "windows":
		nativeStore = client.NewShellProgramFunc("docker-credential-wincred")
	case "darwin":
		nativeStore = client.NewShellProgramFunc("docker-credential-osxkeychain")
	default:
		nativeStore = client.NewShellProgramFunc("docker-credential-secretservice")
	}
	backend.SetCredsLabel("Owl Credentials")
	return HelperStorage{nativeStore}
}

// HelperStorage provides storage for credentials
type HelperStorage struct {
	nativeStore client.ProgramFunc
}

// SetCredentials ...
func (s HelperStorage) SetCredentials(c credentials.Credentials) error {
	normalizedURL, err := NormalizeLdapServerURLWithCred(c.URL(), c.Username())
	if err != nil {
		return err
	}
	credentials := &backend.Credentials{
		ServerURL: normalizedURL,
		Username:  c.Username(),
		Secret:    c.Password(),
	}
	return client.Store(s.nativeStore, credentials)
}

// GetCredentials ...
func (s HelperStorage) GetCredentials(url string, user string) (credentials.Credentials, error) {
	normalizedURL, err := NormalizeLdapServerURLWithCred(url, user)
	if err != nil {
		return nil, err
	}

	storedCreds, err := client.Get(s.nativeStore, normalizedURL)
	if isErrCredentialsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return credentials.NewCredentials(storedCreds.ServerURL, storedCreds.Username, storedCreds.Secret), nil
}

// RemoveCredentials ...
func (s HelperStorage) RemoveCredentials(url string, user string) error {
	normalizedURL, err := NormalizeLdapServerURLWithCred(url, user)
	if err != nil {
		return err
	}

	err = client.Erase(s.nativeStore, normalizedURL)
	if isErrCredentialsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// IsErrCredentialsNotFound returns true if the error
// was caused by not having a set of credentials in a store.
func isErrCredentialsNotFound(err error) bool {
	return backend.IsErrCredentialsNotFound(err)
}
