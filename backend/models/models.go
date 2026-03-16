package models

import "time"

type User struct {
    ID           int       `json:"id"`
    FullName     string    `json:"fullName"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    Phone        *string   `json:"phone,omitempty"`
    Address      *string   `json:"address,omitempty"`
    CreatedAt    time.Time `json:"createdAt"`
}

type Product struct {
    ID          int     `json:"id"`
    Title       string  `json:"title"`
    Price       int     `json:"price"`
    ImageURL    string  `json:"imageUrl"`
    Category    string  `json:"category"`
    Sizes       string  `json:"sizes"`      // строка через запятую
    Material    *string `json:"material,omitempty"`
    Description *string `json:"description,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
}

type Favorite struct {
    UserID    int       `json:"userId"`
    ProductID int       `json:"productId"`
    CreatedAt time.Time `json:"createdAt"`
}

type CartItem struct {
    ID        int       `json:"id"`
    UserID    int       `json:"userId"`
    ProductID int       `json:"productId"`
    Size      string    `json:"size"`
    Quantity  int       `json:"quantity"`
    CreatedAt time.Time `json:"createdAt"`
    // Для ответа с данными товара
    Title     string `json:"title"`
    Price     int    `json:"price"`
    ImageURL  string `json:"imageUrl"`
    Category  string `json:"category"`
}

type Order struct {
    ID         int       `json:"id"`
    UserID     int       `json:"userId"`
    TotalPrice int       `json:"totalPrice"`
    Address    string    `json:"address"`
    CreatedAt  time.Time `json:"createdAt"`
    Items      []OrderItem `json:"items"`
}

type OrderItem struct {
    ID        int       `json:"id"`
    OrderID   int       `json:"orderId"`
    ProductID *int      `json:"productId,omitempty"`
    Title     string    `json:"title"`
    Price     int       `json:"price"`
    Quantity  int       `json:"quantity"`
    Size      *string   `json:"size,omitempty"`
    CreatedAt time.Time `json:"createdAt"`
}