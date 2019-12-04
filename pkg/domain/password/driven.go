package password

// Backend ...
type Backend interface {
	GetPrincipalEmail(userID string) (string, error)
	SetUserPassword(userID string, hashedPassword string) error
}

// MailService ...
type MailService interface {
	SendMail(email string, templateID string, values map[string]string) error
}
