package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/brandon200217/NOTIFY/internal/channel"
	"github.com/brandon200217/NOTIFY/internal/channel/mail"
	"github.com/brandon200217/NOTIFY/internal/channel/mail/factory"
	"github.com/brandon200217/NOTIFY/internal/config"
	"github.com/brandon200217/NOTIFY/internal/handler"
	"github.com/brandon200217/NOTIFY/internal/logger"
	"github.com/brandon200217/NOTIFY/templates"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error cargando configuración: %v", err)
	}

	appLogger := logger.New(logger.Config{
		Level:      cfg.LogLevel,
		Format:     cfg.LogFormat,
		File:       cfg.LogFile,
		MaxSizeMB:  cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAgeDays: cfg.LogMaxAgeDays,
		Compress:   cfg.LogCompress,
	})
	logger.SetDefault(appLogger)

	slog.Info("gomailer iniciando",
		"port", cfg.ServerPort,
		"mail_provider", cfg.MailProvider,
		"log_level", cfg.LogLevel,
		"log_format", cfg.LogFormat)

	engine, err := templates.NewEngine()
	if err != nil {
		slog.Error("error cargando templates", "error", err.Error())
		os.Exit(1)
	}

	mailProvider, err := factory.NewMailProvider(cfg)
	if err != nil {
		slog.Error("error inicializando mail provider", "error", err.Error())
		os.Exit(1)
	}

	mailChannel := mail.NewMailChannel(cfg, mailProvider, engine, cfg.MaxAttachmentSizeMB)
	registry := channel.NewRegistry()
	registry.Register("mail", mailChannel)

	srv := handler.NewServer(cfg, registry)

	slog.Info("server escuchando", "port", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, srv.Router()); err != nil {
		slog.Error("error en server", "error", err.Error())
		os.Exit(1)
	}
}
