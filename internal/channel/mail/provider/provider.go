package provider

import "context"

type Email struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	HTMLBody    string
	Attachments []Attachment
}

type Attachment struct {
	Filename string
	MimeType string
	Content  []byte
}

type MailProvider interface {
	Send(ctx context.Context, email *Email) error
}
