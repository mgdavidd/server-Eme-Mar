package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mgdavidd/server-Eme-Mar/internal/models"
	"github.com/mgdavidd/server-Eme-Mar/internal/services"
	"github.com/mgdavidd/server-Eme-Mar/internal/utils"
)

type ClientHandler struct {
	Service *services.ClientService
}

func NewClientHandler(s *services.ClientService) *ClientHandler {
	return &ClientHandler{Service: s}
}

func (h *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	list, err := h.Service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, "error obteniendo clientes")
		return
	}

	utils.RespondJSON(w, 200, list)
}

func (h *ClientHandler) GetClientById(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.RespondError(w, 400, "id inválido")
		return
	}

	c, err := h.Service.GetById(id)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "cliente no encontrado")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error interno")
		return
	}

	utils.RespondJSON(w, 200, c)
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var c models.Client
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&c); err != nil {
		utils.RespondError(w, 400, "json inválido")
		return
	}

	if err := h.Service.Create(&c); err != nil {
		utils.RespondError(w, 500, "error creando cliente")
		return
	}

	utils.RespondJSON(w, 201, c)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.RespondError(w, 400, "id inválido")
		return
	}

	var c models.Client
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = int64(id)

	err = h.Service.UpdateClient(&c)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "cliente no encontrado")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error actualizando cliente")
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "cliente actualizado"})
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.RespondError(w, 400, "id inválido")
		return
	}

	err = h.Service.DeleteClient(id)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "cliente no encontrado")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error eliminando cliente")
		return
	}

	w.WriteHeader(204)
}
