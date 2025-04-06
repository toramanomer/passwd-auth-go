package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type smtpMailer struct {
	addr     string
	auth     smtp.Auth
	username string
	template *template.Template
	subject  string
	headers  string
}

func newSmptMailer() *smtpMailer {
	var (
		username = os.Getenv("SMTP_USERNAME")
		password = os.Getenv("SMTP_PASSWORD")
		host     = os.Getenv("SMTP_HOST")
		port     = os.Getenv("SMTP_PORT")
	)

	tmpl, err := template.ParseFiles("./core/mailer/template.html")

	if err != nil {
		log.Fatalln(err)
	}

	return &smtpMailer{
		addr:     fmt.Sprintf("%s:%s", host, port),
		auth:     smtp.PlainAuth("", username, password, host),
		username: username,
		template: tmpl,
		subject:  "Subject: Your Verification Code\n",
		headers:  "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
	}
}

func (m *smtpMailer) SendVerificationEmail(to, code string) error {
	var message bytes.Buffer
	if err := m.template.Execute(&message, struct{ Code string }{code}); err != nil {
		return err
	}

	return smtp.SendMail(
		m.addr,
		m.auth,
		m.username,
		[]string{to},
		fmt.Append(nil, m.subject, m.headers, message.String()),
	)
}
