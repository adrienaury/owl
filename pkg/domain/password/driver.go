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

// AssignRandomPassword ...
func (d Driver) AssignRandomPassword(alg string, domain Domain, length uint, userID string) error {
	password, err := d.GetRandomPassword(domain, length)
	if err != nil {
		return err
	}

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
		return fmt.Errorf("invalid password hash algorithm: %v", alg)
	}

	b64 := base64.StdEncoding.EncodeToString(hash)
	result := fmt.Sprintf("{%s}%s", alg, b64)

	emails, err := d.backend.GetUserEmails(userID)
	if err != nil {
		return err
	}

	if len(emails) <= 0 {
		return fmt.Errorf("user has no e-mail, password change is forbidden")
	}

	firstname, err := d.backend.GetUserFirstName(userID)
	if err != nil {
		return err
	}

	lastname, err := d.backend.GetUserLastName(userID)
	if err != nil {
		return err
	}

	if err := d.backend.SetUserPassword(userID, result); err != nil {
		return err
	}

	values := map[string]string{
		"FirstName": firstname,
		"LastName":  lastname,
		"Password":  password,
	}

	errs := []string{}
	for _, email := range emails {
		if err := d.mailService.SendMail(email, "AssignPassword", values); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) == len(emails) {
		return fmt.Errorf(strings.Join(errs, ", "))
	}

	return nil
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
