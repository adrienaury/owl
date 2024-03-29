package password

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Driver is the entry point of the domain that expose methods.
type Driver struct {
	backend Backend
	spusher SecretPusher
}

// NewDriver create a new domain driver with given driven implementations.
func NewDriver(backend Backend, spusher SecretPusher) Driver {
	return Driver{backend, spusher}
}

// AssignRandomPassword to the user with id, the secret is pushed to the user securely.
func (d Driver) AssignRandomPassword(userID string, alg string, domain Domain, length uint) error {
	password, err := d.GetRandomPassword(domain, length)
	if err != nil {
		return err
	}

	if err := d.AssignPassword(userID, alg, password); err != nil {
		return err
	}

	return nil
}

// AssignPassword to the user with id, the secret is pushed to the user securely.
func (d Driver) AssignPassword(userID string, alg string, password string) error {
	mail, err := d.backend.GetPrincipalEmail(userID)
	if err != nil {
		return err
	}

	hash, err := d.GetHash(alg, password)
	if err != nil {
		return err
	}

	if err := d.spusher.CanPushSecret(); err != nil {
		return err
	}

	if err := d.backend.SetUserPassword(userID, hash); err != nil {
		return err
	}

	if err := d.spusher.PushSecret(mail, "user_password", password); err != nil {
		return err
	}

	return nil
}

// GetRandomPassword generate a random password with given length and domain.
func (d Driver) GetRandomPassword(domain Domain, length uint) (string, error) {
	if length == 0 {
		return "", fmt.Errorf("can't generate a 0-length password")
	}
	sb := strings.Builder{}
	for {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(domain.Len())))
		if err != nil {
			return "", err
		}
		sb.WriteRune(domain.Rune(uint(n.Uint64())))
		if sb.Len() == int(length) {
			return sb.String(), nil
		}
	}
}

// GetHashAlgorithms returns list of supported hash algorithms.
func (d Driver) GetHashAlgorithms() []string {
	return []string{
		"SHA224",
		"SSHA224",
		"SHA256",
		"SSHA256",
		"SHA384",
		"SSHA384",
		"SHA512",
		"SSHA512",
		"SHA3-224",
		"SSHA3-224",
		"SHA3-256",
		"SSHA3-256",
		"SHA3-384",
		"SSHA3-384",
		"SHA3-512",
		"SSHA3-512",
		"SHAKE128",
		"SSHAKE128",
		"SHAKE256",
		"SSHAKE256",
	}
}

// GetHash hash the password with the given algorithm.
func (d Driver) GetHash(alg string, password string) (string, error) {
	var hash []byte
	switch alg {
	/* Weak cryptographic primitive blacklisted
	case "MD5":
		hash = makehash(md5.New(), password, false)
	case "SMD5":
		hash = makehash(md5.New(), password, true)
	case "SHA", "SHA1":
		hash = makehash(sha1.New(), password, false)
	case "SSHA", "SSHA1":
		hash = makehash(sha1.New(), password, true) */
	case "SHA224":
		hash = makehash(sha256.New224(), password, false)
	case "SSHA224":
		hash = makehash(sha256.New224(), password, true)
	case "SHA256":
		hash = makehash(sha256.New(), password, false)
	case "SSHA256":
		hash = makehash(sha256.New(), password, true)
	case "SHA384":
		hash = makehash(sha512.New384(), password, false)
	case "SSHA384":
		hash = makehash(sha512.New384(), password, true)
	case "SHA512":
		hash = makehash(sha512.New(), password, false)
	case "SSHA512":
		hash = makehash(sha512.New(), password, true)
	case "SHA3-224":
		hash = makehash(sha3.New224(), password, false)
	case "SSHA3-224":
		hash = makehash(sha3.New224(), password, true)
	case "SHA3-256":
		hash = makehash(sha3.New256(), password, false)
	case "SSHA3-256":
		hash = makehash(sha3.New256(), password, true)
	case "SHA3-384":
		hash = makehash(sha3.New384(), password, false)
	case "SSHA3-384":
		hash = makehash(sha3.New384(), password, true)
	case "SHA3-512":
		hash = makehash(sha3.New512(), password, false)
	case "SSHA3-512":
		hash = makehash(sha3.New512(), password, true)
	case "SHAKE128":
		hash = makeshakehash(sha3.NewShake128(), 32, password, false)
	case "SSHAKE128":
		hash = makeshakehash(sha3.NewShake128(), 32, password, true)
	case "SHAKE256":
		hash = makeshakehash(sha3.NewShake256(), 64, password, false)
	case "SSHAKE256":
		hash = makeshakehash(sha3.NewShake256(), 64, password, true)
	default:
		return "", fmt.Errorf("invalid password hash algorithm: %v", alg)
	}

	b64 := base64.StdEncoding.EncodeToString(hash)
	result := fmt.Sprintf("{%s}%s", alg, b64)

	return result, nil
}

// makehash make a hash of the passphrase with the specified secure hash algorithm
func makeshakehash(alg sha3.ShakeHash, size int, password string, salted bool) []byte {
	h := make([]byte, size)
	_, _ = alg.Write([]byte(password))
	if salted {
		salt := makeSalt(size)
		_, _ = alg.Write(salt)
		_, _ = alg.Read(h)
		return append(h, salt...)
	}
	_, _ = alg.Read(h)
	return h
}

// makehash make a hash of the passphrase with the specified secure hash algorithm
func makehash(alg hash.Hash, password string, salted bool) []byte {
	_, _ = alg.Write([]byte(password))
	if salted {
		salt := makeSalt(alg.Size())
		_, _ = alg.Write(salt)
		h := alg.Sum(nil)
		return append(h, salt...)
	}
	return alg.Sum(nil)
}

// makeSalt make a byte array containing random bytes.
func makeSalt(size int) []byte {
	sbytes := make([]byte, size)
	_, err := rand.Read(sbytes)
	if err != nil {
		panic(err)
	}
	return sbytes
}
