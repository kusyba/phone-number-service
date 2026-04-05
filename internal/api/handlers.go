package api

import (
    "net/http"
    
    "phone-number-service/internal/models"
    "phone-number-service/internal/service"
    
    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
)

type Handlers struct {
    phoneService *service.PhoneService
    validate     *validator.Validate
}

func NewHandlers(phoneService *service.PhoneService) *Handlers {
    return &Handlers{
        phoneService: phoneService,
        validate:     validator.New(),
    }
}

func (h *Handlers) ImportNumbers(c echo.Context) error {
    var req models.ImportRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
    }
    
    if err := h.validate.Struct(req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
    }
    
    stats := h.phoneService.ProcessNumbers(c.Request().Context(), req.Numbers, req.Source)
    
    return c.JSON(http.StatusOK, models.ImportResponse{
        Accepted: stats.Accepted,
        Skipped:  stats.Skipped,
        Errors:   stats.Errors,
    })
}

func (h *Handlers) SearchNumbers(c echo.Context) error {
    var req models.SearchRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid query parameters"})
    }
    
    filters := service.SearchFilters{
        Number:   req.Number,
        Country:  req.Country,
        Region:   req.Region,
        Provider: req.Provider,
        Limit:    req.Limit,
        Offset:   req.Offset,
    }
    
    phones, total, err := h.phoneService.SearchPhones(c.Request().Context(), filters)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    
    return c.JSON(http.StatusOK, models.SearchResponse{
        Data:   phones,
        Total:  total,
        Limit:  filters.Limit,
        Offset: filters.Offset,
    })
}
