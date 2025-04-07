package mailer

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/resend/resend-go/v2"
)

type resendMailer struct {
	client *resend.Client
	tmpl   *template.Template
}

func newResendMailer() *resendMailer {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		panic("RESEND_API_KEY is missing")
	}

	tmpl, err := template.ParseFiles("./core/mailer/template.html")

	if err != nil {
		log.Fatalln(err)
	}

	return &resendMailer{
		client: resend.NewClient(apiKey),
		tmpl:   tmpl,
	}
}

func (mailer *resendMailer) SendVerificationEmail(to, code string) error {
	var message bytes.Buffer
	if err := mailer.tmpl.Execute(&message, struct{ Code string }{code}); err != nil {
		return err
	}

	_, err := mailer.client.Emails.Send(&resend.SendEmailRequest{
		From:    "App <onboarding@resend.dev>",
		To:      []string{to},
		Html:    message.String(),
		Subject: "Login for App",
	})

	return err
}
