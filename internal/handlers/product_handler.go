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

type ProductHandler struct {
	Service *services.ProductService
}

func NewProductHandler(s *services.ProductService) *ProductHandler {
	return &ProductHandler{Service: s}
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	list, err := h.Service.GetAll()
	if err != nil {
		utils.RespondError(w, 500, "Error obteniendo productos")
		return
	}
	utils.RespondJSON(w, 200, list)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var p models.Product

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&p); err != nil {
		utils.RespondError(w, 400, "json inválido")
		return
	}

	if err := validateProduct(p); err != nil {
		utils.RespondError(w, 400, err.Error())
		return
	}

	if err := h.Service.Create(&p); err != nil {
		utils.RespondError(w, 500, "error creando producto")
		return
	}
	utils.RespondJSON(w, 201, p)
}

func validateProduct(p models.Product) error {
	if strings.ReplaceAll(p.Name, " ", "") == "" {
		return errors.New("name invalido")
	}

	if p.Price <= 0 {
		return errors.New("price invalido")
	}

	if len(p.Insumos) == 0 {
		return errors.New("insumos invalidos")
	}

	return nil
}
