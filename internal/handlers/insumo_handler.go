package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mgdavidd/server-Eme-Mar/internal/services"
)

type InsumoHandler struct {
	Service *services.InsumoService
}

func NewInsumoHandler(s *services.InsumoService) *InsumoHandler {
	return &InsumoHandler{
		Service: s,
	}
}

func (h *InsumoHandler) GetAllInsumos(w http.ResponseWriter, r *http.Request) {
	insumos, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "Error Getting Insumos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(insumos)
}
