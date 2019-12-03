package password

// MailService ...
type MailService interface {
	SendMail(email string, templateID string, values map[string]string) error
}
