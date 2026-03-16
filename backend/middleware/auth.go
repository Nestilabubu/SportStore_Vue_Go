package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

const sessionName = "sportshop-session"

func InitSessionStore(key []byte) {
    store = sessions.NewCookieStore(key)
    store.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   86400 * 7, // 7 дней
        HttpOnly: true,
        Secure:   false, // для dev http
        SameSite: http.SameSiteLaxMode,
    }
}

func AuthRequired(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, err := store.Get(r, sessionName)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        userID, ok := session.Values["userID"]
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Добавляем userID в контекст запроса
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetUserID(r *http.Request) int {
    userID, ok := r.Context().Value("userID").(int)
    if !ok {
        return 0
    }
    return userID
}

func GetSessionStore() *sessions.CookieStore {
    return store
}