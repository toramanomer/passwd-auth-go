package mailer

type Mailer interface {
	// Method for sending verification email
	SendVerificationEmail(to, code string) error
}

func NewMailer() Mailer {
	return newSmptMailer()
}
