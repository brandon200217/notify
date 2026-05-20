package discord

import (
	"context"

	"github.com/brandon200217/NOTIFY/internal/config"
	"github.com/brandon200217/NOTIFY/internal/models"
)

type DiscordChannel struct {
	cfg *config.Config
}

func NewDiscordChannel(cfg *config.Config) *DiscordChannel {
	return &DiscordChannel{cfg: cfg}
}

func (m *DiscordChannel) Send(ctx context.Context, req *models.NotifyRequest) error {
	return nil
}
