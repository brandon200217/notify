package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/brandon200217/NOTIFY/internal/models"
)

func (s *Server) handleNotify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decode: %v", err)
		respondError(w, http.StatusBadRequest, "request inválido")
		return
	}

	if req.Token != s.cfg.NotifierToken {
		respondError(w, http.StatusUnauthorized, "token inválido")
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

	if err := ch.Send(context.Background(), &req); err != nil {
		respondError(w, http.StatusInternalServerError, "error al enviar notificación")
		return
	}

	// acá va el dispatch al channel — próximo paso
	respondOK(w, req.Source)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  msg,
		"status": status,
	})
}

func respondOK(w http.ResponseWriter, source string) {
	json.NewEncoder(w).Encode(map[string]string{
		"source": source,
	})
}
