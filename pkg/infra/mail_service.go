package infra

import (
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
	"text/template"
)

// MailService ...
type MailService struct {
	addr string
	from string
}

// NewMailService ...
func NewMailService(addr, from string) MailService {
	return MailService{addr, from}
}

// GetTemplate ...
func (ms MailService) GetTemplate(templateID string) (string, error) {
	dat, err := ioutil.ReadFile(templateID + ".tmpl")
	if os.IsNotExist(err) {
		return "To: {{.to}}\r\n" +
			"Subject: Your password\r\n" +
			"\r\n" +
			"Your password is: {{.password}}\r\n", nil
	} else if err != nil {
		return "", err
	}
	return string(dat), nil
}

// SendMail ...
func (ms MailService) SendMail(email string, templateID string, values map[string]string) error {
	tmpl, err := ms.GetTemplate(templateID)
	if err != nil {
		return err
	}

	sb := strings.Builder{}
	t := template.Must(template.New("mail").Parse(tmpl))
	err = t.Execute(&sb, values)
	if err != nil {
		return err
	}

	err = smtp.SendMail(ms.addr, nil, ms.from, []string{email}, []byte(sb.String()))
	if err != nil {
		return err
	}

	return nil
}
