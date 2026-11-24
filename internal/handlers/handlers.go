package handlers

import (
	"encoding/json"
	"net/http"

	"gold-price-bot.com/v2/internal/database"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello from Go!", "ping": "pong"})
}

func (h *Handler) Recent(w http.ResponseWriter, r *http.Request) {
	latest, err := database.GetLatestGoldPrice()
	if err != nil {
		http.Error(w, "Could not fetch latest price: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latest)
}
