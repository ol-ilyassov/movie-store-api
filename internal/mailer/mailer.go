package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

// Mailer
type Mailer struct {
	dialer *mail.Dialer // SMTP server connection
	sender string       // sender information: "Alice Smith <alice@example.com>"
}

// New
func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second
	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Send
func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// execute defined templates and inject dynamic data.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())      // plain-text body.
	msg.AddAlternative("text/html", htmlBody.String()) // HTML body.
	// AddAlternative() should always be called after SetBody().

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	// retry to send mail 3 times.
	// for i := 1; i <= 3; i++ {
	// 	err = m.dialer.DialAndSend(msg)
	// 	// If everything worked, return nil.
	// 	if nil == err {
	// 		return nil
	// 	}
	// 	// If it didn't work, sleep for a short time and retry.
	// 	time.Sleep(500 * time.Millisecond)
	// }

	return nil
}
