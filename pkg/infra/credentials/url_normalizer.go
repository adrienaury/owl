package credentials

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

var validLdapSchemes = []string{"ldap://", "ldaps://", "tcp://"}

// NormalizeLdapServerURL is opinionated normalization of a string that represents a URL. Returns the URL provided matching the format
// expected when storing a URL in a config. Sets a scheme and port if not present, removes unnecessary trailing
// slashes, etc. Can be used to normalize a URL provided by user input.
func NormalizeLdapServerURL(s string) (string, error) {
	// normalize scheme
	if !hasScheme(s) {
		s = validLdapSchemes[0] + s
	}

	addr, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("not a valid URL: %v", err)
	}

	// normalize host:port
	if strings.Contains(addr.Host, ":") {
		_, port, err := net.SplitHostPort(addr.Host)
		if err != nil {
			return "", fmt.Errorf("not a valid host:port: %v", err)
		}
		_, err = strconv.ParseUint(port, 10, 16)
		if err != nil {
			return "", fmt.Errorf("not a valid port: %v, port numbers must be between 0 and 65535", port)
		}
	} else {
		port := 0
		switch addr.Scheme {
		case "ldap":
			port = 389
		case "ldaps":
			port = 636
		default:
			return "", fmt.Errorf("no port specified")
		}
		addr.Host = net.JoinHostPort(addr.Host, strconv.FormatInt(int64(port), 10))
	}

	// remove trailing slash if that's the only path we have
	if addr.Path == "/" {
		addr.Path = ""
	}

	return addr.String(), nil
}

// NormalizeLdapServerURLWithCred will append username information to the URL
func NormalizeLdapServerURLWithCred(s string, username string) (string, error) {
	s, err := NormalizeLdapServerURL(s)
	if err != nil {
		return s, err
	}

	if username == "" {
		return "", fmt.Errorf("not a valid username: %v", username)
	}

	addr, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("not a valid URL: %v", err)
	}

	addr.User = url.User(username)

	return addr.String(), nil
}

func hasScheme(s string) bool {
	for _, p := range validLdapSchemes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
