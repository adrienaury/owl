package password

// Backend interface.
type Backend interface {
	GetPrincipalEmail(userID string) (string, error)
	SetUserPassword(userID string, hashedPassword string) error
}

// SecretPusher role is to let the user know of a secret in a secure manner.
type SecretPusher interface {
	PushSecret(email string, secretType string, secret string) error
	CanPushSecret() error
}
