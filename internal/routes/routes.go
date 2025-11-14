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
	r.HandleFunc("/clients", clientHandler.GetClients).Methods("GET")
	r.HandleFunc("/clients", clientHandler.CreateClient).Methods("POST")
	r.HandleFunc("/insumos", insumoHandler.GetAllInsumos).Methods("GET")
}
