package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sportshop-backend/db"
	"sportshop-backend/models"
	"sportshop-backend/utils"
	"time"

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

func Register(dbConn *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        var userID int
        err = dbConn.QueryRow(`
            INSERT INTO users (full_name, email, password_hash, phone, address)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id`,
            req.FullName, req.Email, string(hashed), req.Phone, req.Address,
        ).Scan(&userID)

        if err != nil {
            http.Error(w, "Email already exists", http.StatusConflict)
            return
        }

        token, err := utils.GenerateSessionToken()
        if err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        expiresAt := time.Now().Add(15 * time.Minute)
        if err := db.CreateSession(dbConn, token, userID, expiresAt); err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:     "session_token",
            Value:    token,
            Path:     "/",
            HttpOnly: true,
            Secure:   true,                     // обязательно для HTTPS
            SameSite: http.SameSiteNoneMode,   // разрешить кросс-домен
            MaxAge:   86400,
        })

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]int{"id": userID})
    }
}

func Login(dbConn *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        var user models.User
        err := dbConn.QueryRow(`
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

        token, err := utils.GenerateSessionToken()
        if err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        expiresAt := time.Now().Add(15 * time.Minute)
        if err := db.CreateSession(dbConn, token, user.ID, expiresAt); err != nil {
            http.Error(w, "Server error", http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:     "session_token",
            Value:    token,
            Path:     "/",
            HttpOnly: true,
            Secure:   true,                     // обязательно для HTTPS
            SameSite: http.SameSiteNoneMode,   // разрешить кросс-домен
            MaxAge:   86400,
        })

        json.NewEncoder(w).Encode(user)
    }
}

func Logout(dbConn *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("session_token")
        if err == nil {
            db.DeleteSessionByToken(dbConn, cookie.Value)
        }
        http.SetCookie(w, &http.Cookie{
            Name:     "session_token",
            Value:    "",
            Path:     "/",
            HttpOnly: true,
            MaxAge:   -1,
        })
        w.WriteHeader(http.StatusOK)
    }
}