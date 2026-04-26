package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
    // Пробуем получить DATABASE_URL от Railway
    databaseURL := os.Getenv("DATABASE_URL")
    psqlInfo := ""
    
    if databaseURL != "" {
        // Railway предоставляет полный URL
        psqlInfo = databaseURL
        log.Println("Using DATABASE_URL from Railway")
    } else {
        // Fallback для локальной разработки
        host := os.Getenv("DB_HOST")
        if host == "" {
            host = "localhost"
        }
        port := os.Getenv("DB_PORT")
        if port == "" {
            port = "5432"
        }
        user := os.Getenv("DB_USER")
        if user == "" {
            user = "postgres"
        }
        password := os.Getenv("DB_PASSWORD")
        dbname := os.Getenv("DB_NAME")
        if dbname == "" {
            dbname = "railway"
        }
        
        psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
            host, port, user, password, dbname)
        log.Println("Using local database connection")
    }
    
    log.Printf("Connecting to database with connection string: %s", maskPassword(psqlInfo))
    
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Error opening database: ", err)
    }
    
    // Настройка пула соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Проверяем подключение с повторными попытками
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        err = db.Ping()
        if err == nil {
            break
        }
        log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
        if i == maxRetries-1 {
            log.Fatal("Error connecting to database after retries: ", err)
        }
        time.Sleep(3 * time.Second)
    }
    
    log.Println("Successfully connected to database")
    return db
}

// Вспомогательная функция для скрытия пароля в логах
func maskPassword(connStr string) string {
    // Ищем password= в строке подключения
    if strings.Contains(connStr, "password=") {
        parts := strings.Split(connStr, "password=")
        if len(parts) > 1 {
            passwordPart := strings.Split(parts[1], " ")[0]
            masked := strings.Repeat("*", len(passwordPart))
            return parts[0] + "password=" + masked + strings.TrimPrefix(parts[1], passwordPart)
        }
    }
    return connStr
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