package main

import (
	"log"
	"net/http"
	"os"

	"sportshop-backend/db"
	"sportshop-backend/handlers"
	"sportshop-backend/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
    database := db.InitDB()
    defer database.Close()

    r := mux.NewRouter()

    r.HandleFunc("/api/register", handlers.Register(database)).Methods("POST")
    r.HandleFunc("/api/login", handlers.Login(database)).Methods("POST")
    r.HandleFunc("/api/logout", handlers.Logout(database)).Methods("POST")
    r.HandleFunc("/api/products", handlers.GetProducts(database)).Methods("GET")
    r.HandleFunc("/api/products/{id}", handlers.GetProduct(database)).Methods("GET")

    auth := r.PathPrefix("/api").Subrouter()
    auth.Use(middleware.AuthRequired(database))
    auth.HandleFunc("/profile", handlers.GetProfile(database)).Methods("GET")
    auth.HandleFunc("/profile", handlers.UpdateProfile(database)).Methods("PUT")
    auth.HandleFunc("/favorites", handlers.GetFavorites(database)).Methods("GET")
    auth.HandleFunc("/favorites", handlers.AddFavorite(database)).Methods("POST")
    auth.HandleFunc("/favorites/{productId}", handlers.RemoveFavorite(database)).Methods("DELETE")
    auth.HandleFunc("/cart", handlers.GetCart(database)).Methods("GET")
    auth.HandleFunc("/cart", handlers.AddToCart(database)).Methods("POST")
    auth.HandleFunc("/cart/{itemId}", handlers.UpdateCartItem(database)).Methods("PUT")
    auth.HandleFunc("/cart/{itemId}", handlers.RemoveFromCart(database)).Methods("DELETE")
    auth.HandleFunc("/orders", handlers.GetOrders(database)).Methods("GET")
    auth.HandleFunc("/orders", handlers.CreateOrder(database)).Methods("POST")
    auth.HandleFunc("/statistics", handlers.GetStatistics(database)).Methods("GET")

    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"},
        AllowCredentials: true,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
    })

    handler := c.Handler(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server started on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}