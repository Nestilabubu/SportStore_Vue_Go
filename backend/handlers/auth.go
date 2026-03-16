package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sportshop-backend/middleware"
	"sportshop-backend/models"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
    FullName string  `json:"fullName"`
    Email    string  `json:"email"`
    Password string  `json:"password"`
    Phone    *string `json:"phone"`
    Address  *string `json:"address"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func Register(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        // Хеширование пароля
        hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        var userID int
        err = db.QueryRow(`
            INSERT INTO users (full_name, email, password_hash, phone, address)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id`,
            req.FullName, req.Email, string(hashed), req.Phone, req.Address,
        ).Scan(&userID)

        if err != nil {
            http.Error(w, "Email already exists", http.StatusConflict)
            return
        }

        // Устанавливаем сессию
		session, _ := middleware.GetSessionStore().Get(r, "sportshop-session")	
        session.Values["userID"] = userID
        session.Save(r, w)

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]int{"id": userID})
    }
}

func Login(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        var user models.User
        err := db.QueryRow(`
            SELECT id, full_name, email, password_hash, phone, address, created_at
            FROM users WHERE email = $1`, req.Email,
        ).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Phone, &user.Address, &user.CreatedAt)

        if err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        session, _ := middleware.GetSessionStore().Get(r, "sportshop-session")
        session.Values["userID"] = user.ID
        session.Save(r, w)

        json.NewEncoder(w).Encode(user)
    }
}

func Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := middleware.GetSessionStore().Get(r, "sportshop-session")
    session.Values["userID"] = nil
    session.Options.MaxAge = -1 // удалить
    session.Save(r, w)
    w.WriteHeader(http.StatusOK)
}