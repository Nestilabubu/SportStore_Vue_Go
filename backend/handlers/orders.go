package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sportshop-backend/middleware"
	"sportshop-backend/models"
	"time"
)

func GetOrders(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        rows, err := db.Query(`
            SELECT id, total_price, address, created_at
            FROM orders
            WHERE user_id = $1
            ORDER BY created_at DESC`, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        orders := []models.Order{}
        for rows.Next() {
            var o models.Order
            o.UserID = userID
            err := rows.Scan(&o.ID, &o.TotalPrice, &o.Address, &o.CreatedAt)
            if err != nil {
                continue
            }

            // Загружаем элементы заказа с image_url
            itemRows, err := db.Query(`
                SELECT id, product_id, title, price, quantity, size, image_url, created_at
                FROM order_items
                WHERE order_id = $1`, o.ID)
            if err == nil {
                items := []models.OrderItem{}
                for itemRows.Next() {
                    var oi models.OrderItem
                    oi.OrderID = o.ID
                    err := itemRows.Scan(&oi.ID, &oi.ProductID, &oi.Title, &oi.Price, &oi.Quantity, &oi.Size, &oi.ImageUrl, &oi.CreatedAt)
                    if err != nil {
                        continue
                    }
                    items = append(items, oi)
                }
                o.Items = items
                itemRows.Close()
            }

            orders = append(orders, o)
        }

        json.NewEncoder(w).Encode(orders)
    }
}

func CreateOrder(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        // Получаем корзину пользователя с image_url
        rows, err := db.Query(`
            SELECT c.product_id, c.size, c.quantity, p.title, p.price, p.image_url
            FROM cart c
            JOIN products p ON c.product_id = p.id
            WHERE c.user_id = $1`, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        type cartItem struct {
            ProductID int
            Size      string
            Quantity  int
            Title     string
            Price     int
            ImageUrl  string
        }
        items := []cartItem{}
        totalPrice := 0
        for rows.Next() {
            var ci cartItem
            err := rows.Scan(&ci.ProductID, &ci.Size, &ci.Quantity, &ci.Title, &ci.Price, &ci.ImageUrl)
            if err != nil {
                continue
            }
            items = append(items, ci)
            totalPrice += ci.Price * ci.Quantity
        }

        if len(items) == 0 {
            http.Error(w, "Cart is empty", http.StatusBadRequest)
            return
        }

        // Получаем адрес пользователя
        var address string
        err = db.QueryRow(`SELECT address FROM users WHERE id = $1`, userID).Scan(&address)
        if err != nil {
            address = "Не указан"
        }

        // Создаём заказ
        tx, err := db.Begin()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        var orderID int
        err = tx.QueryRow(`
            INSERT INTO orders (user_id, total_price, address, created_at)
            VALUES ($1, $2, $3, $4)
            RETURNING id`, userID, totalPrice, address, time.Now(),
        ).Scan(&orderID)
        if err != nil {
            tx.Rollback()
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        for _, item := range items {
            _, err = tx.Exec(`
                INSERT INTO order_items (order_id, product_id, title, price, quantity, size, image_url, created_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
                orderID, item.ProductID, item.Title, item.Price, item.Quantity, item.Size, item.ImageUrl, time.Now())
            if err != nil {
                tx.Rollback()
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        // Очищаем корзину
        _, err = tx.Exec(`DELETE FROM cart WHERE user_id = $1`, userID)
        if err != nil {
            tx.Rollback()
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        err = tx.Commit()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]int{"orderId": orderID})
    }
}