package handler

import (
	"encoding/json"
	"net/http"
)

func respondError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  msg,
		"status": status,
	})
}

func respondOK(w http.ResponseWriter, source string) {
	_ = json.NewEncoder(w).Encode(map[string]string{
		"source": source,
	})
}
