package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    
    "phone-number-service/internal/config"
    
    _ "github.com/lib/pq"
)

// NewDBConnection создает соединение с PostgreSQL
func NewDBConnection(cfg *config.Config) (*sql.DB, error) {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    for i := 0; i < 30; i++ {
        if err := db.Ping(); err == nil {
            log.Println("Database ping successful")
            return db, nil
        }
        log.Printf("Waiting for database... (%d/30)", i+1)
        time.Sleep(2 * time.Second)
    }
    
    return nil, fmt.Errorf("failed to ping database after 30 attempts")
}

// Queries - структура для работы с sqlc
type Queries struct {
    db *sql.DB
}

// New - создает новый экземпляр Queries
func New(db *sql.DB) *Queries {
    return &Queries{db: db}
}
