package api

import (
    "net/http"
    
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(e *echo.Echo, handlers *Handlers) {
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())
    
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
    })
    
    api := e.Group("/api")
    numbers := api.Group("/numbers")
    numbers.POST("/import", handlers.ImportNumbers)
    numbers.GET("/search", handlers.SearchNumbers)
}
