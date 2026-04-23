package api

import (
    "net/http"
    
    "github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, handlers *Handlers, groupHandlers *GroupHandlers) {
    // Health check
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
    })
    
    api := e.Group("/api")
    
    // Phone endpoints
    numbers := api.Group("/numbers")
    numbers.POST("/import", handlers.ImportNumbers)
    numbers.GET("/search", handlers.SearchNumbers)
    
    // Format endpoints
    phones := api.Group("/phones")
    phones.GET("/:id/format", handlers.FormatPhoneByID)
    phones.POST("/format", handlers.FormatPhoneByValue)
    
    // Group endpoints
    groups := api.Group("/groups")
    groups.POST("", groupHandlers.CreateGroup)
    groups.GET("", groupHandlers.GetGroups)
    groups.GET("/:id", groupHandlers.GetGroupByID)
    groups.PUT("/:id", groupHandlers.UpdateGroup)
    groups.DELETE("/:id", groupHandlers.DeleteGroup)
    groups.POST("/:id/users/:userId", groupHandlers.AddUserToGroup)
    groups.DELETE("/:id/users/:userId", groupHandlers.RemoveUserFromGroup)
    
    // User endpoints
    users := api.Group("/users")
    users.GET("/:userId/groups", groupHandlers.GetUserGroups)
    
    // Me endpoint
    api.GET("/me", handlers.GetMe)
}
