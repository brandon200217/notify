package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brandon200217/NOTIFY/internal/channel"
	"github.com/brandon200217/NOTIFY/internal/channel/mail"
	"github.com/brandon200217/NOTIFY/internal/channel/mail/factory"
	"github.com/brandon200217/NOTIFY/internal/config"
	"github.com/brandon200217/NOTIFY/internal/handler"
	"github.com/brandon200217/NOTIFY/internal/logger"
	"github.com/brandon200217/NOTIFY/internal/middleware"
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

	rate := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	srv := handler.NewServer(cfg, registry, rate)

	httpServer := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: srv.Router(),
	}

	go func() {
		slog.Info("server escuchando", "port", cfg.ServerPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error en server", "error", err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	slog.Info("señal recibida, iniciando shutdown",
		"signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("error durante shutdown", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("servidor cerrado correctamente")

}
