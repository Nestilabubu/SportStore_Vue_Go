package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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

    // Публичные маршруты
    r.HandleFunc("/api/register", handlers.Register(database)).Methods("POST")
    r.HandleFunc("/api/login", handlers.Login(database)).Methods("POST")
    r.HandleFunc("/api/logout", handlers.Logout(database)).Methods("POST")
    r.HandleFunc("/api/products", handlers.GetProducts(database)).Methods("GET")
    r.HandleFunc("/api/products/{id}", handlers.GetProduct(database)).Methods("GET")

    // Защищённые маршруты
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

    // Настройка CORS для продакшена
    allowedOrigins := []string{
        "https://sport-store-vue-go.vercel.app",
        "https://sport-store-vue-go-git-main-nestilabubus-projects.vercel.app",
        "http://localhost:5173",
        "http://localhost:3000",
    }
    
    // Добавляем возможность указывать CORS через переменную окружения
    if corsOrigins := os.Getenv("CORS_ORIGINS"); corsOrigins != "" {
        allowedOrigins = strings.Split(corsOrigins, ",")
    }
    
    // Настройка CORS с динамическим разрешением origin
    c := cors.New(cors.Options{
    AllowOriginFunc: func(origin string) bool {
        // Разрешаем локальную разработку
        if origin == "http://localhost:5173" || origin == "http://localhost:3000" {
            return true
        }
        // Основной production домен
        if origin == "https://sport-store-vue-go.vercel.app" {
            return true
        }
        // Все превью-домены Vercel (заканчиваются на .vercel.app)
        if strings.HasSuffix(origin, ".vercel.app") {
            return true
        }
        return false
    },
    AllowCredentials: true,
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Content-Type", "Authorization", "Cookie"},
    ExposedHeaders:   []string{"Set-Cookie"},
})

    handler := c.Handler(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server started on port %s", port)
    log.Printf("CORS allowed origins: %v", allowedOrigins)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}