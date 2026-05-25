package models

import "fmt"

type NotifyRequest struct {
	Source string `json:"source"`
	Type   string `json:"type"`

	// Mail
	Receivers   []string     `json:"receivers,omitempty"`
	Cc          []string     `json:"cc,omitempty"`
	Bcc         []string     `json:"bcc,omitempty"`
	Subject     string       `json:"subject,omitempty"`
	Body        string       `json:"body,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`

	// Camino 2 — template de gomailer
	TemplateID   string                 `json:"template_id,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`

	// Slack / Discord
	Text    string `json:"text,omitempty"`
	Channel string `json:"channel,omitempty"`
	Url     string `json:"url,omitempty"`
}

type Attachment struct {
	Filename      string `json:"filename"`
	MimeType      string `json:"mime_type"`
	Base64Content string `json:"base64_content"`
}

func (r *NotifyRequest) Validate() error {
	if r.Source == "" {
		return fmt.Errorf("el campo source es requerido")
	}
	if r.Type == "" {
		return fmt.Errorf("el campo type es requerido")
	}

	switch r.Type {
	case "mail":
		return r.validateMail()
	case "slack", "discord":
		return r.validateMessaging()
	default:
		return fmt.Errorf("type inválido: %s — valores aceptados: mail, slack, discord", r.Type)
	}
}

func (r *NotifyRequest) validateMail() error {
	if len(r.Receivers) == 0 {
		return fmt.Errorf("mail requiere al menos un receiver")
	}
	if r.Subject == "" {
		return fmt.Errorf("mail requiere subject")
	}
	if r.Body == "" && r.TemplateID == "" {
		return fmt.Errorf("mail requiere body o template_id")
	}
	return nil
}

func (r *NotifyRequest) validateMessaging() error {
	if r.Text == "" {
		return fmt.Errorf("%s requiere text", r.Type)
	}
	if r.Channel == "" {
		return fmt.Errorf("%s requiere channel", r.Type)
	}
	return nil
}
