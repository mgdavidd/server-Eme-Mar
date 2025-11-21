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

func (h *MoveHandler) GetAllMoves(w http.ResponseWriter, r *http.Request) {
	data, err := h.Service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, "error obteniendo movimientos")
		return
	}

	utils.RespondJSON(w, 200, data)
}

func (h *MoveHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	data, err := h.Service.GetBalance()
	if err != nil {
		utils.RespondError(w, 500, "error obteniendo Saldo/Fiado")
		return
	}

	utils.RespondJSON(w, 200, data)
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

func (h *MoveHandler) Sell(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var sale models.Sale

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&sale); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			utils.RespondError(w, 400, "campo desconocido en el cuerpo de la solicitud")
			return
		}
		utils.RespondError(w, 400, "error al parsear el cuerpo de la solicitud")
		return
	}

	if sale.ClientId <= 0 {
		utils.RespondError(w, 400, "client_id inválido")
		return
	}

	if len(sale.Items) == 0 {
		utils.RespondError(w, 400, "debes enviar al menos 1 producto")
		return
	}

	for _, item := range sale.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			utils.RespondError(w, 400, "producto inválido en items")
			return
		}
	}

	err := h.Service.Sell(sale)

	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			utils.RespondError(w, 404, "cliente o producto no encontrado")
			return
		}
		utils.RespondError(w, 500, "error procesando venta")
		return
	}

	utils.RespondJSON(w, 200, map[string]string{
		"message": "venta procesada correctamente",
	})
}

func (h *MoveHandler) PayCredit(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req struct {
		CreditSaleID int64   `json:"credit_sale_id"`
		Amount       float64 `json:"amount"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			utils.RespondError(w, 400, "campo desconocido en el cuerpo de la solicitud")
			return
		}
		utils.RespondError(w, 400, "error al parsear el cuerpo de la solicitud")
		return
	}

	if req.CreditSaleID <= 0 {
		utils.RespondError(w, 400, "credit_sale_id inválido")
		return
	}
	if req.Amount <= 0 {
		utils.RespondError(w, 400, "amount inválido")
		return
	}

	err := h.Service.PayCredit(req.CreditSaleID, req.Amount)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			utils.RespondError(w, 404, "venta a crédito no encontrada")
			return
		}
		if errors.Is(err, services.ErrInvalidInput) {
			utils.RespondError(w, 400, "monto inválido o supera el saldo pendiente")
			return
		}
		utils.RespondError(w, 500, "error procesando abono")
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"message": "abono procesado correctamente"})
}
