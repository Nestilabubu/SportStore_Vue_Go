package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"sportshop-backend/middleware"
	"sportshop-backend/models"

	"github.com/gorilla/mux"
)

func GetFavorites(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        rows, err := db.Query(`
            SELECT p.id, p.title, p.price, p.image_url, p.category, p.sizes, p.material, p.description
            FROM favorites f
            JOIN products p ON f.product_id = p.id
            WHERE f.user_id = $1
            ORDER BY f.created_at DESC`, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        favorites := []models.Product{}
        for rows.Next() {
            var p models.Product
            err := rows.Scan(&p.ID, &p.Title, &p.Price, &p.ImageURL, &p.Category, &p.Sizes, &p.Material, &p.Description)
            if err != nil {
                continue
            }
            favorites = append(favorites, p)
        }

        json.NewEncoder(w).Encode(favorites)
    }
}

type AddFavoriteRequest struct {
    ProductID int `json:"productId"`
}

func AddFavorite(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        var req AddFavoriteRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        _, err := db.Exec(`
            INSERT INTO favorites (user_id, product_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING`, userID, req.ProductID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    }
}

func RemoveFavorite(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)
        vars := mux.Vars(r)
        productId, _ := strconv.Atoi(vars["productId"])

        _, err := db.Exec(`DELETE FROM favorites WHERE user_id = $1 AND product_id = $2`, userID, productId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}