package handler

import (
	"context"
	"errors"

	"github.com/brandon200217/NOTIFY/internal/channel"
	"github.com/brandon200217/NOTIFY/internal/config"
	"github.com/brandon200217/NOTIFY/internal/middleware"
	"github.com/brandon200217/NOTIFY/internal/models"
)

const testToken = "modelovistacontrolador"

type mockChannel struct {
	shouldFail bool
}

func (m *mockChannel) Send(ctx context.Context, req *models.NotifyRequest) error {
	if m.shouldFail {
		return errors.New("error simulado de envío")
	}
	return nil
}

func newTestServer(mailChannel *mockChannel) *Server {
	cfg := &config.Config{
		NotifierToken:       testToken,
		MaxRequestSizeMB:    10,
		MaxAttachmentSizeMB: 5,
	}

	registry := channel.NewRegistry()
	registry.Register("mail", mailChannel)

	rateLimiter := middleware.NewRateLimiter(100, 200)

	return NewServer(cfg, registry, rateLimiter)
}

func authHeader() string {
	return "Bearer " + testToken
}
