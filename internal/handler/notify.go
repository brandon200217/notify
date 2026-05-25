package handler

import (
	"crypto/subtle"
	"encoding/json"
	"log"
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
		respondError(w, http.StatusUnauthorized, "header Authorization requerido")
		return
	}

	provided := []byte(token)
	expected := []byte(s.cfg.NotifierToken)

	if subtle.ConstantTimeCompare(provided, expected) != 1 {
		respondError(w, http.StatusBadRequest, "request inválido")
	}

	var req models.NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decode: %v", err)
		respondError(w, http.StatusBadRequest, "request inválido")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	ch, err := s.registry.Get(req.Type)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ch.Send(r.Context(), &req); err != nil {
		respondError(w, http.StatusInternalServerError, "error al enviar notificación")
		return
	}

	respondOK(w, req.Source)
}
