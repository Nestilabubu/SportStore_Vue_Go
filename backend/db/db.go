package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    if host == "" {
        host = "localhost"
    }
    if port == "" {
        port = "5432"
    }
    if user == "" {
        user = "sportuser"
    }
    if password == "" {
        password = "sportpass"
    }
    if dbname == "" {
        dbname = "sportshop"
    }

    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Error opening database: ", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Error connecting to database: ", err)
    }

    log.Println("Successfully connected to database")
    return db
}

func CreateSession(db *sql.DB, token string, userID int, expiresAt time.Time) error {
    _, err := db.Exec("INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)", token, userID, expiresAt)
    return err
}

func GetSessionByToken(db *sql.DB, token string) (sessionID int, userID int, expiresAt time.Time, err error) {
    err = db.QueryRow("SELECT id, user_id, expires_at FROM sessions WHERE token = $1 AND expires_at > $2", token, time.Now()).
        Scan(&sessionID, &userID, &expiresAt)
    if err != nil {
        return 0, 0, time.Time{}, err
    }
    return
}

func UpdateSessionExpiry(db *sql.DB, token string, newExpiry time.Time) error {
    _, err := db.Exec("UPDATE sessions SET expires_at = $1 WHERE token = $2", newExpiry, token)
    return err
}

func DeleteSessionByToken(db *sql.DB, token string) error {
    _, err := db.Exec("DELETE FROM sessions WHERE token = $1", token)
    return err
}