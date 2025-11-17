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

type InsumoHandler struct {
	Service *services.InsumoService
}

func NewInsumoHandler(s *services.InsumoService) *InsumoHandler {
	return &InsumoHandler{Service: s}
}

func (h *InsumoHandler) GetAllInsumos(w http.ResponseWriter, r *http.Request) {
	data, err := h.Service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, "error obteniendo insumos")
		return
	}

	utils.RespondJSON(w, 200, data)
}

func (h *InsumoHandler) GetByIdInsumos(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.RespondError(w, 400, "id inválido")
		return
	}

	insumo, err := h.Service.GetById(id)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "insumo no encontrado")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error interno")
		return
	}

	utils.RespondJSON(w, 200, insumo)
}

func (h *InsumoHandler) CreateInsumo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var in models.Insumo

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&in); err != nil {
		utils.RespondError(w, 400, "json inválido")
		return
	}

	if err := validateInsumo(in); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	if err := h.Service.Create(&in); err != nil {
		utils.RespondError(w, 500, "error creando insumo")
		return
	}

	utils.RespondJSON(w, 201, in)
}

func (h *InsumoHandler) UpdateInsumo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.RespondError(w, 400, "id inválido")
		return
	}

	var in models.Insumo

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&in); err != nil {
		utils.RespondError(w, 400, "json inválido")
		return
	}

	if err := validateInsumo(in); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	in.ID = int64(id)
	err = h.Service.Update(&in)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "insumo no encontrado")
		return
	}

	if err != nil {
		utils.RespondError(w, 500, "Error actualizando insumo")
		return
	}
	utils.RespondJSON(w, 200, in)

}

// VALIDACIÓN
func validateInsumo(i models.Insumo) error {
	if i.Name == "" {
		return errors.New("el campo 'name' es obligatorio")
	}
	if i.Um == "" {
		return errors.New("el campo 'um' es obligatorio")
	}
	if i.Stock < 0 {
		return errors.New("el campo 'stock' no puede ser negativo")
	}
	if i.MinStock < 0 {
		return errors.New("el campo 'min_stock' no puede ser negativo")
	}
	if i.UnitPrice <= 0 {
		return errors.New("el campo 'unit_price' debe ser mayor a 0")
	}
	return nil
}
