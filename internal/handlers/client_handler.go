package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mgdavidd/server-Eme-Mar/internal/models"
	"github.com/mgdavidd/server-Eme-Mar/internal/services"
)

type ClientHandler struct {
	Service *services.ClientService
}

func NewClientHandler(s *services.ClientService) *ClientHandler {
	return &ClientHandler{Service: s}
}

func (h *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "Error getting clients", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(clients)
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var c models.Client
	json.NewDecoder(r.Body).Decode(&c)

	err := h.Service.Create(&c)
	if err != nil {
		http.Error(w, "Error creating client", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(c)
}
