package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"sportshop-backend/db"
	"time"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthRequired(dbConn *sql.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            cookie, err := r.Cookie("session_token")
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            _, userID, _, err := db.GetSessionByToken(dbConn, cookie.Value)
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            newExpiry := time.Now().Add(15 * time.Minute)
            db.UpdateSessionExpiry(dbConn, cookie.Value, newExpiry)

            http.SetCookie(w, &http.Cookie{
                Name:     "session_token",
                Value:    cookie.Value,
                Path:     "/",
                HttpOnly: true,
                Secure:   true,                     // обязательно для HTTPS
                SameSite: http.SameSiteNoneMode,   // разрешить кросс-домен
                MaxAge:   86400,
            })

            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserID(r *http.Request) int {
    userID, ok := r.Context().Value(UserIDKey).(int)
    if !ok {
        return 0
    }
    return userID
}