package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/mgdavidd/server-Eme-Mar/internal/db"
	"github.com/mgdavidd/server-Eme-Mar/internal/handlers"
	"github.com/mgdavidd/server-Eme-Mar/internal/routes"
	"github.com/mgdavidd/server-Eme-Mar/internal/services"
)

func main() {
	fmt.Println("Starting server...")

	// DB
	database := db.ConnectDB()
	db.RunMigrations(database)

	// Services
	clientService := services.NewClientService(database)
	insumoService := services.NewInsumoService(database)
	moveService := services.NewMoveService(database)
	productService := services.NewProductService(database)

	// Handlers
	clientHandler := handlers.NewClientHandler(clientService)
	insumoHandler := handlers.NewInsumoHandler(insumoService)
	moveHandler := handlers.NewMoveHandler(moveService)
	productHandler := handlers.NewProductHandler(productService)

	// Router
	r := mux.NewRouter()
	routes.RegisterRoutes(r, clientHandler, insumoHandler, moveHandler, productHandler)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false, // IMPORTANTE: con "*" no puedes usar credenciales
	})

	handler := c.Handler(r)

	// Start server
	log.Println("Server running on http://localhost:3000")
	http.ListenAndServe(":3000", handler)
}
