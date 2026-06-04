package factory

import (
	"fmt"

	"github.com/brandon200217/NOTIFY/internal/channel/mail/provider"
	"github.com/brandon200217/NOTIFY/internal/channel/mail/provider/smtp"
	"github.com/brandon200217/NOTIFY/internal/config"
)

func NewMailProvider(cfg *config.Config) (provider.MailProvider, error) {
	switch cfg.MailProvider {
	case "smtp":
		return smtp.New(smtp.Config{
			Host:              cfg.SMTPHost,
			Port:              cfg.SMTPPort,
			Username:          cfg.SMTPUsername,
			Password:          cfg.SMTPPassword,
			From:              cfg.SMTPFrom,
			ConnectTimeoutSec: cfg.SMTPConnectTimeoutSec,
			SendTimeoutSec:    cfg.SMTPSendTimeoutSec,
		}), nil

	default:
		return nil, fmt.Errorf("mail provider no soportado: %q (válidos: smtp)", cfg.MailProvider)
	}
}
