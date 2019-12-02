package password

// Backend ...
type Backend interface {
	SetUserPassword(userID string, hashedPassword string) error
	GetUserEmails(id string) ([]string, error)
	GetUserFirstName(id string) (string, error)
	GetUserLastName(id string) (string, error)
}

// MailService ...
type MailService interface {
	SendMail(email string, templateID string, values map[string]string) error
}
