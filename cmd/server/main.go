package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "phone-number-service/internal/api"
    "phone-number-service/internal/config"
    "phone-number-service/internal/database"
    "phone-number-service/internal/service"
    "phone-number-service/pkg/logger"
    
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
    
    cfg := config.Load()
    logger.InitLogger(cfg.LogLevel)
    
    db, err := database.NewDBConnection(cfg)
    if err != nil {
        logger.Global.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    logger.Global.Info("Database connected successfully")
    
    phoneService := service.NewPhoneService(db)
    handlers := api.NewHandlers(phoneService)
    
    e := echo.New()
    api.SetupRoutes(e, handlers)
    
    go func() {
        if err := e.Start(":" + cfg.Port); err != nil {
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
