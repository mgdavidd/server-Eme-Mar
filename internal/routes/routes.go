package routes

import (
	"github.com/gorilla/mux"
	"github.com/mgdavidd/server-Eme-Mar/internal/handlers"
)

func RegisterRoutes(
	r *mux.Router,
	clientHandler *handlers.ClientHandler,
	insumoHandler *handlers.InsumoHandler,
) {

	// --- CLIENTES ---
	clientRoutes := r.PathPrefix("/clients").Subrouter()
	clientRoutes.HandleFunc("", clientHandler.GetClients).Methods("GET")
	clientRoutes.HandleFunc("", clientHandler.CreateClient).Methods("POST")
	clientRoutes.HandleFunc("/{id}", clientHandler.UpdateClient).Methods("PUT")
	clientRoutes.HandleFunc("/{id}", clientHandler.DeleteClient).Methods("DELETE")
	clientRoutes.HandleFunc("/{id}", clientHandler.GetClientById).Methods("GET")

	// --- INSUMOS ---
	insumoRoutes := r.PathPrefix("/insumos").Subrouter()
	insumoRoutes.HandleFunc("", insumoHandler.GetAllInsumos).Methods("GET")
	insumoRoutes.HandleFunc("", insumoHandler.CreateInsumo).Methods("POST")
	insumoRoutes.HandleFunc("/{id}", insumoHandler.GetByIdInsumos).Methods("GET")
	insumoRoutes.HandleFunc("/{id}", insumoHandler.UpdateInsumo).Methods("PUT")

}
