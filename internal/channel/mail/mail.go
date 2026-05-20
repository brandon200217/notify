package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/brandon200217/NOTIFY/internal/config"
	"github.com/brandon200217/NOTIFY/internal/models"
	"github.com/brandon200217/NOTIFY/templates"
	mail "github.com/xhit/go-simple-mail/v2"
)

type MailChannel struct {
	cfg    *config.Config
	engine *templates.Engine
}

func NewMailChannel(cfg *config.Config, engine *templates.Engine) *MailChannel {
	return &MailChannel{cfg: cfg, engine: engine}
}

func (m *MailChannel) Send(ctx context.Context, req *models.NotifyRequest) error {
	client, err := m.newClient()
	if err != nil {
		return fmt.Errorf("error conectando al servidor SMTP: %w", err)
	}

	smtpClient, err := client.Connect()
	if err != nil {
		return fmt.Errorf("error estableciendo conexión SMTP: %w", err)
	}
	defer smtpClient.Close()

	email, err := m.buildEmail(req)
	if err != nil {
		return fmt.Errorf("error construyendo email: %w", err)
	}
	if err := email.Send(smtpClient); err != nil {
		return fmt.Errorf("error enviando email: %w", err)
	}

	return nil
}

func (m *MailChannel) newClient() (*mail.SMTPServer, error) {
	server := mail.NewSMTPClient()
	server.Host = m.cfg.SMTPHost
	server.Port = m.cfg.SMTPPort
	server.Username = m.cfg.SMTPUsername
	server.Password = m.cfg.SMTPPassword
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	return server, nil
}

func (m *MailChannel) buildEmail(req *models.NotifyRequest) (*mail.Email, error) {
	email := mail.NewMSG()

	email.SetFrom(m.cfg.SMTPFrom)
	email.AddTo(req.Receivers...)
	email.SetSubject(req.Subject)
	email.SetBody(mail.TextHTML, req.Body)

	body := req.Body
	if req.TemplateID != "" {
		rendered, err := m.engine.Render(req.TemplateID, req.TemplateData)
		if err != nil {
			return nil, err
		}
		body = rendered
	}

	email.SetBody(mail.TextHTML, body)

	// CC y BCC opcionales
	if len(req.Cc) > 0 {
		email.AddCc(req.Cc...)
	}
	if len(req.Bcc) > 0 {
		email.AddBcc(req.Bcc...)
	}

	// adjuntos opcionales
	for _, att := range req.Attachments {
		decoded, err := base64.StdEncoding.DecodeString(att.Base64Content)
		if err != nil {
			continue
		}
		email.Attach(&mail.File{
			Name:     att.Filename,
			MimeType: att.MimeType,
			Data:     decoded,
			Inline:   false,
		})
	}

	return email, nil
}
