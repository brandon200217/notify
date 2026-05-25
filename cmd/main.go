package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brandon200217/NOTIFY/internal/channel"
	"github.com/brandon200217/NOTIFY/internal/channel/mail"
	"github.com/brandon200217/NOTIFY/internal/channel/mail/factory"
	"github.com/brandon200217/NOTIFY/internal/config"
	"github.com/brandon200217/NOTIFY/internal/handler"
	"github.com/brandon200217/NOTIFY/templates"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error cargando configuración: %v", err)
	}
	engine, err := templates.NewEngine()
	if err != nil {
		log.Fatalf("error cargando templates: %v", err)
	}

	mailProvider, err := factory.NewMailProvider(cfg)
	if err != nil {
		log.Fatalf("error inicializando mail provider: %v", err)
	}

	mailChannel := mail.NewMailChannel(cfg, mailProvider, engine, cfg.MaxAttachmentSizeMB)
	registry := channel.NewRegistry()
	registry.Register("mail", mailChannel)

	srv := handler.NewServer(cfg, registry)

	fmt.Printf("gomailer corriendo en %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, srv.Router()))
}
