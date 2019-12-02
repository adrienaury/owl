package infra

import "fmt"

// MailService ...
type MailService struct {
}

// NewMailService ...
func NewMailService() MailService {
	return MailService{}
}

// SendMail ...
func (ms MailService) SendMail(email string, templateID string, values map[string]string) error {
	fmt.Println(values)
	return nil
}
