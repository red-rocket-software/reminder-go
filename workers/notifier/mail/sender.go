package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
	smtpAuthAddress   string
	smtpServerAddress string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword, smtpAuthAddress, smtpServerAddress string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
		smtpAuthAddress:   smtpAuthAddress,
		smtpServerAddress: smtpServerAddress,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := &email.Email{
		From:    fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress),
		To:      to,
		Bcc:     bcc,
		Cc:      cc,
		Subject: subject,
		HTML:    []byte(content),
	}

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, sender.smtpAuthAddress)
	t := &tls.Config{InsecureSkipVerify: true}
	return e.SendWithStartTLS(sender.smtpServerAddress, smtpAuth, t)
}
