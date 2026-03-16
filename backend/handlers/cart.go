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

type AddToCartRequest struct {
    ProductID int    `json:"productId"`
    Size      string `json:"size"`
    Quantity  int    `json:"quantity"` // опционально, по умолчанию 1
}

func GetCart(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        rows, err := db.Query(`
            SELECT c.id, c.product_id, c.size, c.quantity, c.created_at,
                   p.title, p.price, p.image_url, p.category
            FROM cart c
            JOIN products p ON c.product_id = p.id
            WHERE c.user_id = $1
            ORDER BY c.created_at DESC`, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        cart := []models.CartItem{}
        for rows.Next() {
            var ci models.CartItem
            err := rows.Scan(&ci.ID, &ci.ProductID, &ci.Size, &ci.Quantity, &ci.CreatedAt,
                &ci.Title, &ci.Price, &ci.ImageURL, &ci.Category)
            if err != nil {
                continue
            }
            ci.UserID = userID
            cart = append(cart, ci)
        }

        json.NewEncoder(w).Encode(cart)
    }
}

func AddToCart(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        var req AddToCartRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        if req.Quantity <= 0 {
            req.Quantity = 1
        }

        // Проверяем, есть ли уже такой товар с таким размером
        var existingID int
        err := db.QueryRow(`
            SELECT id FROM cart
            WHERE user_id = $1 AND product_id = $2 AND size = $3`,
            userID, req.ProductID, req.Size).Scan(&existingID)

        if err == nil {
            // Обновляем количество
            _, err = db.Exec(`
                UPDATE cart SET quantity = quantity + $1
                WHERE id = $2`, req.Quantity, existingID)
        } else {
            // Вставляем новую запись
            _, err = db.Exec(`
                INSERT INTO cart (user_id, product_id, size, quantity)
                VALUES ($1, $2, $3, $4)`,
                userID, req.ProductID, req.Size, req.Quantity)
        }

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}

func UpdateCartItem(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)
        vars := mux.Vars(r)
        itemId, _ := strconv.Atoi(vars["itemId"])

        var req struct {
            Quantity int `json:"quantity"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        if req.Quantity <= 0 {
            // Удаляем, если количество 0
            db.Exec(`DELETE FROM cart WHERE id = $1 AND user_id = $2`, itemId, userID)
        } else {
            _, err := db.Exec(`
                UPDATE cart SET quantity = $1
                WHERE id = $2 AND user_id = $3`, req.Quantity, itemId, userID)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        w.WriteHeader(http.StatusOK)
    }
}

func RemoveFromCart(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)
        vars := mux.Vars(r)
        itemId, _ := strconv.Atoi(vars["itemId"])

        _, err := db.Exec(`DELETE FROM cart WHERE id = $1 AND user_id = $2`, itemId, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}