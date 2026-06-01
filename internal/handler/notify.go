package handler

import (
	"crypto/subtle"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brandon200217/NOTIFY/internal/models"
	"github.com/brandon200217/NOTIFY/utilities"
)

func (s *Server) handleNotify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	maxBytes := int64(s.cfg.MaxRequestSizeMB) << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	token := utilities.ExtractBearerToken(r)
	if token == "" {
		slog.WarnContext(r.Context(), "header Authorization faltante",
			"remote_addr", r.RemoteAddr)
		respondError(w, http.StatusUnauthorized, "header Authorization requerido")
		return
	}

	provided := []byte(token)
	expected := []byte(s.cfg.NotifierToken)

	if subtle.ConstantTimeCompare(provided, expected) != 1 {
		slog.WarnContext(r.Context(), "token inválido",
			"remote_addr", r.RemoteAddr)
		respondError(w, http.StatusUnauthorized, "token inválido")
		return
	}

	var req models.NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(r.Context(), "error decodificando request",
			"error", err.Error(),
			"remote_addr", r.RemoteAddr)
		respondError(w, http.StatusBadRequest, "request inválido")
		return
	}

	if err := req.Validate(); err != nil {
		slog.WarnContext(r.Context(), "request inválido",
			"source", req.Source,
			"type", req.Type,
			"error", err.Error())
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ch, err := s.registry.Get(req.Type)
	if err != nil {
		slog.WarnContext(r.Context(), "canal no soportado",
			"type", req.Type,
			"source", req.Source)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ch.Send(r.Context(), &req); err != nil {
		slog.ErrorContext(r.Context(), "error enviando notificación",
			"source", req.Source,
			"type", req.Type,
			"error", err.Error())
		respondError(w, http.StatusInternalServerError, "error al enviar notificación")
		return
	}

	slog.InfoContext(r.Context(), "notificación procesada",
		"source", req.Source,
		"type", req.Type,
		"receivers_count", len(req.Receivers),
		"template_id", req.TemplateID)

	respondOK(w, req.Source)
}
