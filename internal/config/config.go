package config

import (
	"fmt"

	"github.com/brandon200217/NOTIFY/utilities"
)

type Config struct {
	// Server
	ServerPort string
	SMTPFrom   string

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	// Slack
	SlackWebhookURL string

	// Discord
	DiscordChannelsPath string

	// Auth
	NotifierToken string

	// Logging
	LogPath  string
	LogLevel string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:          utilities.GetEnvOrDefault("SERVER_PORT", ":8080"),
		SMTPHost:            utilities.GetEnvOrDefault("SMTP_HOST", ""),
		SMTPPort:            utilities.GetEnvOrDefaultInt("SMTP_PORT", 587),
		SMTPUsername:        utilities.GetEnvOrDefault("SMTP_USERNAME", ""),
		SMTPPassword:        utilities.GetEnvOrDefault("SMTP_PASSWORD", ""),
		SlackWebhookURL:     utilities.GetEnvOrDefault("SLACK_WEBHOOK_URL", ""),
		DiscordChannelsPath: utilities.GetEnvOrDefault("DISCORD_CHANNELS_PATH", "./discord_channels.json"),
		NotifierToken:       utilities.GetEnvOrDefault("NOTIFIER_TOKEN", ""),
		LogPath:             utilities.GetEnvOrDefault("LOG_PATH", "./"),
		LogLevel:            utilities.GetEnvOrDefault("LOG_LEVEL", "ERROR"),
		SMTPFrom:            utilities.GetEnvOrDefault("SMTP_FROM", ""),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	required := map[string]string{
		"NOTIFIER_TOKEN": c.NotifierToken,
		"SMTP_HOST":      c.SMTPHost,
		"SMTP_USERNAME":  c.SMTPUsername,
		"SMTP_PASSWORD":  c.SMTPPassword,
		"SMTP_FROM":      c.SMTPFrom,
	}

	for key, val := range required {
		if val == "" {
			return fmt.Errorf("variable de entorno requerida faltante: %s", key)
		}
	}

	return nil
}
