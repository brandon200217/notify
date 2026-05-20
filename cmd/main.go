package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brandon200217/NOTIFY/internal/channel"
	"github.com/brandon200217/NOTIFY/internal/channel/discord"
	"github.com/brandon200217/NOTIFY/internal/channel/mail"
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
	registry := channel.NewRegistry()
	registry.Register("mail", mail.NewMailChannel(cfg, engine))
	registry.Register("discord", discord.NewDiscordChannel(cfg))
	srv := handler.NewServer(cfg, registry)

	fmt.Printf("gomailer corriendo en %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, srv.Router()))
}
