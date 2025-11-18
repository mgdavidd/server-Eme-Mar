package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mgdavidd/server-Eme-Mar/internal/models"
	"github.com/mgdavidd/server-Eme-Mar/internal/services"
	"github.com/mgdavidd/server-Eme-Mar/internal/utils"
)

type MoveHandler struct {
	Service *services.MovementService
}

func NewMoveHandler(s *services.MovementService) *MoveHandler {
	return &MoveHandler{Service: s}
}

func (h *MoveHandler) Supply(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var supply models.Supply

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&supply); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			utils.RespondError(w, 400, "campo desconocido en el cuerpo de la solicitud")
			return
		}
		utils.RespondError(w, 400, "error al parsear el cuerpo de la solicitud")
		return
	}

	err := h.Service.Supply(supply)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			utils.RespondError(w, 400, "entrada inválida")
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			utils.RespondError(w, 404, "insumo no encontrado")
			return
		}
		utils.RespondError(w, 500, "error interno del servidor")
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "surtido realizado con éxito"})
}
