package password

// Backend ...
type Backend interface {
	GetPrincipalEmail(userID string) (string, error)
	SetUserPassword(userID string, hashedPassword string) error
}

// SecretPusher ...
type SecretPusher interface {
	PushSecret(email string, secretType string, secret string) error
}
