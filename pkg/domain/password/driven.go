package password

// Backend ...
type Backend interface {
	SetEmailVerificationURI(userID string, secret string, expire int64) error
	VerifyEmail(userID string, secret string) error
	GetVerifiedEmail(userID string) (string, error)
	SetUserPassword(userID string, hashedPassword string) error
}

// MailService ...
type MailService interface {
	SendMail(email string, templateID string, values map[string]string) error
}
