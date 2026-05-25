package smtp

import (
	"context"
	"fmt"
	"time"

	"github.com/brandon200217/NOTIFY/internal/channel/mail/provider"
	mail "github.com/xhit/go-simple-mail/v2"
)

type SMTPProvider struct {
	host     string
	port     int
	username string
	password string
	from     string
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func New(cfg Config) *SMTPProvider {
	return &SMTPProvider{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		from:     cfg.From,
	}
}

func (p *SMTPProvider) Send(ctx context.Context, email *provider.Email) error {
	server := mail.NewSMTPClient()
	server.Host = p.host
	server.Port = p.port
	server.Username = p.username
	server.Password = p.password
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return fmt.Errorf("error conectando a SMTP: %w", err)
	}
	defer smtpClient.Close()

	msg := p.buildMessage(email)

	if err := msg.Send(smtpClient); err != nil {
		return fmt.Errorf("error enviando email por SMTP: %w", err)
	}

	return nil
}

func (p *SMTPProvider) buildMessage(email *provider.Email) *mail.Email {
	msg := mail.NewMSG()
	msg.SetFrom(p.from)
	msg.AddTo(email.To...)
	msg.SetSubject(email.Subject)
	msg.SetBody(mail.TextHTML, email.HTMLBody)

	if len(email.Cc) > 0 {
		msg.AddCc(email.Cc...)
	}
	if len(email.Bcc) > 0 {
		msg.AddBcc(email.Bcc...)
	}

	for _, att := range email.Attachments {
		msg.Attach(&mail.File{
			Name:     att.Filename,
			MimeType: att.MimeType,
			Data:     att.Content,
			Inline:   false,
		})
	}

	return msg
}
