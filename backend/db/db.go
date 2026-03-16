package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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