package password

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Driver ...
type Driver struct {
	backend     Backend
	mailService MailService
}

// NewDriver ...
func NewDriver(backend Backend, mailService MailService) Driver {
	return Driver{backend, mailService}
}

// AssignRandomPassword ...
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

// AssignPassword ...
func (d Driver) AssignPassword(userID string, alg string, password string) error {
	mail, err := d.backend.GetVerifiedEmail(userID)
	if err != nil {
		return err
	}

	hash, err := d.GetHash(alg, password)
	if err != nil {
		return err
	}

	if err := d.backend.SetUserPassword(userID, hash); err != nil {
		return err
	}

	values := map[string]string{
		"password": password,
		"to":       mail,
	}

	if err := d.mailService.SendMail(mail, "AssignPassword", values); err != nil {
		return err
	}

	return nil
}

// GetRandomPassword ...
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

// GetHash ...
func (d Driver) GetHash(alg string, password string) (string, error) {
	var hash []byte
	switch alg {
	case "MD5":
		hash = makehash(md5.New(), password, false)
	case "SMD5":
		hash = makehash(md5.New(), password, true)
	case "SHA", "SHA1":
		hash = makehash(sha1.New(), password, false)
	case "SSHA", "SSHA1":
		hash = makehash(sha1.New(), password, true)
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
	default:
		return "", fmt.Errorf("invalid password hash algorithm: %v", alg)
	}

	b64 := base64.StdEncoding.EncodeToString(hash)
	result := fmt.Sprintf("{%s}%s", alg, b64)

	return result, nil
}

// makehash make a hash of the passphrase with the specified secure hash algorithm
func makehash(alg hash.Hash, password string, salted bool) []byte {
	alg.Write([]byte(password))
	if salted {
		salt := makeSalt()
		alg.Write(salt)
		h := alg.Sum(nil)
		return append(h, salt...)
	}
	return alg.Sum(nil)
}

// makeSalt make a 4 byte array containing random bytes.
func makeSalt() []byte {
	sbytes := make([]byte, 4)
	rand.Read(sbytes)
	return sbytes
}
