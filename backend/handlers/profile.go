package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sportshop-backend/middleware"
	"sportshop-backend/models"
)

func GetProfile(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        var user models.User
        err := db.QueryRow(`
            SELECT id, full_name, email, phone, address, created_at
            FROM users WHERE id = $1`, userID,
        ).Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.Address, &user.CreatedAt)

        if err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(user)
    }
}

type UpdateProfileRequest struct {
    FullName *string `json:"fullName"`
    Phone    *string `json:"phone"`
    Address  *string `json:"address"`
}

func UpdateProfile(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        var req UpdateProfileRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        // Динамическое обновление только переданных полей
        if req.FullName != nil {
            db.Exec("UPDATE users SET full_name = $1 WHERE id = $2", *req.FullName, userID)
        }
        if req.Phone != nil {
            db.Exec("UPDATE users SET phone = $1 WHERE id = $2", *req.Phone, userID)
        }
        if req.Address != nil {
            db.Exec("UPDATE users SET address = $1 WHERE id = $2", *req.Address, userID)
        }

        w.WriteHeader(http.StatusOK)
    }
}