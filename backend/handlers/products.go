package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"sportshop-backend/models"

	"github.com/gorilla/mux"
)

func GetProducts(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Парсинг query параметров
        sortBy := r.URL.Query().Get("sortBy")
        category := r.URL.Query().Get("category")
        search := r.URL.Query().Get("search") // title LIKE
        minPrice := r.URL.Query().Get("minPrice")
        maxPrice := r.URL.Query().Get("maxPrice")
        size := r.URL.Query().Get("size")

        query := "SELECT id, title, price, image_url, category, sizes, material, description, created_at FROM products WHERE 1=1"
        args := []interface{}{}
        argId := 1

        if category != "" && category != "all" {
            query += " AND category = $" + strconv.Itoa(argId)
            args = append(args, category)
            argId++
        }

        if search != "" {
            query += " AND title ILIKE '%' || $" + strconv.Itoa(argId) + " || '%'"
            args = append(args, search)
            argId++
        }

        if minPrice != "" {
            min, _ := strconv.Atoi(minPrice)
            query += " AND price >= $" + strconv.Itoa(argId)
            args = append(args, min)
            argId++
        }

        if maxPrice != "" {
            max, _ := strconv.Atoi(maxPrice)
            query += " AND price <= $" + strconv.Itoa(argId)
            args = append(args, max)
            argId++
        }

        if size != "" {
            query += " AND sizes LIKE '%' || $" + strconv.Itoa(argId) + " || '%'"
            args = append(args, size)
            argId++
        }

        // Сортировка
        switch sortBy {
        case "price":
            query += " ORDER BY price ASC"
        case "-price":
            query += " ORDER BY price DESC"
        case "title":
            query += " ORDER BY title ASC"
        default:
            query += " ORDER BY title ASC"
        }

        rows, err := db.Query(query, args...)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        products := []models.Product{}
        for rows.Next() {
            var p models.Product
            err := rows.Scan(&p.ID, &p.Title, &p.Price, &p.ImageURL, &p.Category, &p.Sizes, &p.Material, &p.Description, &p.CreatedAt)
            if err != nil {
                continue
            }
            products = append(products, p)
        }

        json.NewEncoder(w).Encode(products)
    }
}

func GetProduct(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, _ := strconv.Atoi(vars["id"])

        var p models.Product
        err := db.QueryRow(`
            SELECT id, title, price, image_url, category, sizes, material, description, created_at
            FROM products WHERE id = $1`, id,
        ).Scan(&p.ID, &p.Title, &p.Price, &p.ImageURL, &p.Category, &p.Sizes, &p.Material, &p.Description, &p.CreatedAt)

        if err != nil {
            http.Error(w, "Product not found", http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(p)
    }
}