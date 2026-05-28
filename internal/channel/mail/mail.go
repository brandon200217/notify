package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/brandon200217/NOTIFY/internal/channel/mail/provider"
	"github.com/brandon200217/NOTIFY/internal/config"
	"github.com/brandon200217/NOTIFY/internal/models"
	"github.com/brandon200217/NOTIFY/templates"
	"github.com/gabriel-vasile/mimetype"
)

type MailChannel struct {
	provider           provider.MailProvider
	engine             *templates.Engine
	maxAttachmentBytes int
}

func NewMailChannel(cfg *config.Config, p provider.MailProvider, engine *templates.Engine, maxAttachmentMB int) *MailChannel {
	return &MailChannel{
		provider:           p,
		engine:             engine,
		maxAttachmentBytes: maxAttachmentMB << 20,
	}
}

func (m *MailChannel) Send(ctx context.Context, req *models.NotifyRequest) error {
	email, err := m.buildEmail(req)
	if err != nil {
		return fmt.Errorf("error construyendo email: %w", err)
	}

	return m.provider.Send(ctx, email)
}

func (m *MailChannel) buildEmail(req *models.NotifyRequest) (*provider.Email, error) {
	body, err := m.resolveBody(req)
	if err != nil {
		return nil, err
	}

	email := &provider.Email{
		To:       req.Receivers,
		Cc:       req.Cc,
		Bcc:      req.Bcc,
		Subject:  req.Subject,
		HTMLBody: body,
	}

	for _, att := range req.Attachments {
		decoded, err := base64.StdEncoding.DecodeString(att.Base64Content)
		if err != nil {
			return nil, fmt.Errorf("error decodificando adjunto %s: %w", att.Filename, err)
		}

		if len(decoded) > m.maxAttachmentBytes {
			return nil, fmt.Errorf("Attachments supera 10mb")
		}

		mtype := mimetype.Detect(decoded)
		if !mimeMatches(att.MimeType, mtype.String()) {
			return nil, fmt.Errorf("mime type declarado (%s) no coincide con el contenido real (%s) en %s",
				att.MimeType, mtype.String(), att.Filename)
		}

		email.Attachments = append(email.Attachments, provider.Attachment{
			Filename: att.Filename,
			MimeType: att.MimeType,
			Content:  decoded,
		})
	}

	return email, nil
}

func (m *MailChannel) resolveBody(req *models.NotifyRequest) (string, error) {
	if req.TemplateID != "" {
		return m.engine.Render(req.TemplateID, req.TemplateData)
	}
	return req.Body, nil
}

func mimeMatches(declared, detected string) bool {
	declaredBase := stripMimeParams(declared)
	detectedBase := stripMimeParams(detected)

	if detectedBase == "application/octet-stream" {
		return true
	}

	return declaredBase == detectedBase
}

func stripMimeParams(mimeType string) string {
	if idx := strings.Index(mimeType, ";"); idx != -1 {
		return strings.TrimSpace(mimeType[:idx])
	}
	return strings.TrimSpace(mimeType)
}
