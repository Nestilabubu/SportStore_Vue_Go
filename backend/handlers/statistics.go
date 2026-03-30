package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sportshop-backend/middleware"
)

type Statistics struct {
    TotalOrders      int            `json:"totalOrders"`
    TotalSpent       int            `json:"totalSpent"`
    AverageOrder     float64        `json:"averageOrder"`
    FavoriteCategory *string        `json:"favoriteCategory"`
    TotalItemsBought int            `json:"totalItemsBought"`
    TotalUniqueItems int            `json:"totalUniqueItems"`
    CategoryStats    map[string]struct {
        Count      int `json:"count"`
        TotalSpent int `json:"totalSpent"`
    } `json:"categoryStats"`
}

func GetStatistics(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := middleware.GetUserID(r)

        var stats Statistics
        stats.CategoryStats = make(map[string]struct {
            Count      int `json:"count"`
            TotalSpent int `json:"totalSpent"`
        })

        err := db.QueryRow(`
            SELECT COUNT(*), COALESCE(SUM(total_price), 0)
            FROM orders
            WHERE user_id = $1`, userID).Scan(&stats.TotalOrders, &stats.TotalSpent)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if stats.TotalOrders > 0 {
            stats.AverageOrder = float64(stats.TotalSpent) / float64(stats.TotalOrders)
        }

        err = db.QueryRow(`
            SELECT COALESCE(SUM(quantity), 0), COUNT(DISTINCT product_id)
            FROM order_items oi
            JOIN orders o ON oi.order_id = o.id
            WHERE o.user_id = $1`, userID).Scan(&stats.TotalItemsBought, &stats.TotalUniqueItems)
        if err != nil {
        }

        rows, err := db.Query(`
            SELECT p.category, SUM(oi.quantity) as cnt, SUM(oi.price * oi.quantity) as spent
            FROM order_items oi
            JOIN orders o ON oi.order_id = o.id
            JOIN products p ON oi.product_id = p.id
            WHERE o.user_id = $1
            GROUP BY p.category`, userID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var maxCategory string
        maxCount := 0
        for rows.Next() {
            var cat string
            var cnt, spent int
            rows.Scan(&cat, &cnt, &spent)
            stats.CategoryStats[cat] = struct {
                Count      int `json:"count"`
                TotalSpent int `json:"totalSpent"`
            }{cnt, spent}
            if cnt > maxCount {
                maxCount = cnt
                maxCategory = cat
            }
        }

        if maxCount > 0 {
            stats.FavoriteCategory = &maxCategory
        }

        json.NewEncoder(w).Encode(stats)
    }
}