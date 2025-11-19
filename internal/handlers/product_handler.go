package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

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

func (h *ProductHandler) GetByIdProducts(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.RespondError(w, 400, "id inválido")
		return
	}

	product, err := h.Service.GetById(id)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "producto no encontrado")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error interno")
		return
	}

	utils.RespondJSON(w, 200, product)
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

// nombre, precio, foto
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	pidStr := vars["id"]
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		utils.RespondError(w, 400, "product id inválido")
		return
	}
	var p models.ProductSimple
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		utils.RespondError(w, 400, "json inválido")
		return
	}
	p.ID = int64(pid)

	if strings.ReplaceAll(p.Name, " ", "") == "" {
		utils.RespondError(w, 400, "name inválido")
		return
	}
	if p.Price <= 0 {
		utils.RespondError(w, 400, "price inválido")
		return
	}
	if err := h.Service.Update(p); err != nil {
		utils.RespondError(w, 500, "error actualizando producto")
		return
	}
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.RespondError(w, 400, "id inválido")
		return
	}

	err = h.Service.Delete(id)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "producto no encontrado")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error eliminando producto")
		return
	}

	w.WriteHeader(204)
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

// UpdateProductInsumo handles updating the cantidad_insumo for a product-insumo relation
func (h *ProductHandler) UpdateProductInsumo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	pidStr := vars["id"]
	iidStr := vars["insumo_id"]

	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		utils.RespondError(w, 400, "product id inválido")
		return
	}
	iid, err := strconv.Atoi(iidStr)
	if err != nil || iid <= 0 {
		utils.RespondError(w, 400, "insumo id inválido")
		return
	}

	var body struct {
		Quantity float64 `json:"quantity"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		utils.RespondError(w, 400, "json inválido")
		return
	}

	if body.Quantity < 0 {
		utils.RespondError(w, 400, "quantity inválida")
		return
	}

	// This endpoint updates quantity only; fails if relation doesn't exist.
	err = h.Service.UpdateInsumoQuantity(int64(pid), int64(iid), body.Quantity)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "relación no encontrada")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error actualizando cantidad")
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"status": "ok"})
}

// AddProductInsumo handles creating a producto_insumos relation or updating it if exists
// POST /products/{id}/insumos/{insumo_id}
func (h *ProductHandler) AddProductInsumo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	pidStr := vars["id"]
	iidStr := vars["insumo_id"]

	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		utils.RespondError(w, 400, "product id inválido")
		return
	}
	iid, err := strconv.Atoi(iidStr)
	if err != nil || iid <= 0 {
		utils.RespondError(w, 400, "insumo id inválido")
		return
	}

	var body struct {
		Quantity float64 `json:"quantity"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		utils.RespondError(w, 400, "json inválido")
		return
	}

	if body.Quantity < 0 {
		utils.RespondError(w, 400, "quantity inválida")
		return
	}

	err = h.Service.UpdateOrCreateInsumo(int64(pid), int64(iid), body.Quantity)
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "producto o insumo no encontrado")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error creando o actualizando relación")
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"status": "ok"})
}

func (h *ProductHandler) DeleteProductInsumo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pidStr := vars["id"]
	iidStr := vars["insumo_id"]

	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid <= 0 {
		utils.RespondError(w, 400, "product id inválido")
		return
	}
	iid, err := strconv.Atoi(iidStr)
	if err != nil || iid <= 0 {
		utils.RespondError(w, 400, "insumo id inválido")
		return
	}

	err = h.Service.RemoveInsumoFromProduct(int64(pid), int64(iid))
	if errors.Is(err, services.ErrNotFound) {
		utils.RespondError(w, 404, "relación no encontrada")
		return
	}
	if err != nil {
		utils.RespondError(w, 500, "error eliminando relación")
		return
	}

	utils.RespondJSON(w, 200, map[string]string{"status": "ok"})
}
