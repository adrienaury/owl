package credentials

import (
	"runtime"

	"github.com/adrienaury/owl/pkg/helpers/normalizer"
	"github.com/docker/docker-credential-helpers/client"
	"github.com/docker/docker-credential-helpers/credentials"
)

var nativeStore client.ProgramFunc

func init() {
	switch runtime.GOOS {
	case "windows":
		nativeStore = client.NewShellProgramFunc("docker-credential-wincred")
	case "darwin":
		nativeStore = client.NewShellProgramFunc("docker-credential-osxkeychain")
	default:
		nativeStore = client.NewShellProgramFunc("docker-credential-secretservice")
	}
	credentials.SetCredsLabel("Owl Credentials")
}

// SetCredentials store credentials in a secure vault
func SetCredentials(url string, user string, password string) error {
	normalizedURL, err := normalizer.NormalizeLdapServerURLWithCred(url, user)
	if err != nil {
		return err
	}
	c := &credentials.Credentials{
		ServerURL: normalizedURL,
		Username:  user,
		Secret:    password,
	}
	return client.Store(nativeStore, c)
}

// GetCredentials retrieve credentials from a secure vault
func GetCredentials(url string, user string) (*credentials.Credentials, error) {
	normalizedURL, err := normalizer.NormalizeLdapServerURLWithCred(url, user)
	if err != nil {
		return nil, err
	}

	storedCreds, err := client.Get(nativeStore, normalizedURL)
	if err != nil {
		return nil, err
	}
	return storedCreds, nil
}

// IsErrCredentialsNotFound returns true if the error
// was caused by not having a set of credentials in a store.
func IsErrCredentialsNotFound(err error) bool {
	return credentials.IsErrCredentialsNotFound(err)
}
