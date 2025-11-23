package routes

import (
	"github.com/gorilla/mux"
	"github.com/mgdavidd/server-Eme-Mar/internal/handlers"
)

func RegisterRoutes(
	r *mux.Router,
	clientHandler *handlers.ClientHandler,
	insumoHandler *handlers.InsumoHandler,
	movesHandler *handlers.MoveHandler,
	productHandler *handlers.ProductHandler,
) {

	// --- CLIENTES ---
	clientRoutes := r.PathPrefix("/clients").Subrouter()
	clientRoutes.HandleFunc("", clientHandler.GetClients).Methods("GET")
	clientRoutes.HandleFunc("", clientHandler.CreateClient).Methods("POST")
	clientRoutes.HandleFunc("/debt", clientHandler.GetIndebtedClient).Methods("GET")
	clientRoutes.HandleFunc("/{id}", clientHandler.UpdateClient).Methods("PUT")
	clientRoutes.HandleFunc("/{id}", clientHandler.DeleteClient).Methods("DELETE")
	clientRoutes.HandleFunc("/{id}", clientHandler.GetClientById).Methods("GET")

	// --- INSUMOS ---
	insumoRoutes := r.PathPrefix("/insumos").Subrouter()
	insumoRoutes.HandleFunc("", insumoHandler.GetAllInsumos).Methods("GET")
	insumoRoutes.HandleFunc("", insumoHandler.CreateInsumo).Methods("POST")
	insumoRoutes.HandleFunc("/{id}", insumoHandler.GetByIdInsumos).Methods("GET")
	insumoRoutes.HandleFunc("/{id}", insumoHandler.UpdateInsumo).Methods("PUT")
	insumoRoutes.HandleFunc("/{id}", insumoHandler.DeleteInsumo).Methods("DELETE")

	// --- PRODUCTOS ---
	productRoutes := r.PathPrefix("/products").Subrouter()
	productRoutes.HandleFunc("", productHandler.CreateProduct).Methods("POST")
	productRoutes.HandleFunc("", productHandler.GetAllProducts).Methods("GET")
	productRoutes.HandleFunc("/{id}", productHandler.GetByIdProducts).Methods("GET")
	productRoutes.HandleFunc("/{id}", productHandler.UpdateProduct).Methods("PUT")
	productRoutes.HandleFunc("/{id}", productHandler.DeleteProduct).Methods("DELETE")
	productRoutes.HandleFunc("/{id}/insumos/{insumo_id}", productHandler.AddProductInsumo).Methods("POST")
	productRoutes.HandleFunc("/{id}/insumos/{insumo_id}", productHandler.UpdateProductInsumo).Methods("PUT")
	productRoutes.HandleFunc("/{id}/insumos/{insumo_id}", productHandler.DeleteProductInsumo).Methods("DELETE")

	// --- MOVIMIENTOS ---
	movesRoutes := r.PathPrefix("/moves").Subrouter()
	movesRoutes.HandleFunc("", movesHandler.Supply).Methods("POST")
	movesRoutes.HandleFunc("/sell", movesHandler.Sell).Methods("POST")
	movesRoutes.HandleFunc("/pay/credit", movesHandler.PayCredit).Methods("POST")
	movesRoutes.HandleFunc("", movesHandler.GetAllMoves).Methods("GET")
	movesRoutes.HandleFunc("/account", movesHandler.GetAccount).Methods("GET")
	movesRoutes.HandleFunc("/client/{id}", movesHandler.GetClientCreditSales).Methods("GET")
	movesRoutes.HandleFunc("/credit/payments/{sale_id}", movesHandler.GetCreditPayments).Methods("GET")

}
