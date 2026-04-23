package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "phone-number-service/internal/api"
    "phone-number-service/internal/config"
    "phone-number-service/internal/service"
    "phone-number-service/pkg/logger"
    
    "github.com/go-playground/validator/v10"
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    _ "github.com/lib/pq"
)

type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
    
    cfg := config.Load()
    logger.InitLogger(cfg.LogLevel)
    
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        logger.Global.Fatal("Failed to connect to database:", err)
    }
    
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    for i := 0; i < 30; i++ {
        if err := db.Ping(); err == nil {
            log.Println("Database ping successful")
            break
        }
        if i == 29 {
            logger.Global.Fatal("Failed to ping database after 30 attempts:", err)
        }
        log.Printf("Waiting for database... (%d/30)", i+1)
        time.Sleep(2 * time.Second)
    }
    defer db.Close()
    
    logger.Global.Info("Database connected successfully")
    
    phoneService := service.NewPhoneService(db)
    groupService := service.NewGroupService(db)
    
    handlers := api.NewHandlers(phoneService, groupService)
    groupHandlers := api.NewGroupHandlers(groupService)
    
    e := echo.New()
    e.HideBanner = true
    e.Validator = &CustomValidator{validator: validator.New()}
    
    e.Use(middleware.Recover())
    e.Use(middleware.Logger())
    e.Use(middleware.CORS())
    
    // Передаём оба хендлера
    api.SetupRoutes(e, handlers, groupHandlers)
    
    go func() {
        if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
            logger.Global.Infof("Server stopped: %v", err)
        }
    }()
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    logger.Global.Info("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := e.Shutdown(ctx); err != nil {
        logger.Global.Fatal("Server forced to shutdown:", err)
    }
    
    logger.Global.Info("Server exited")
}
