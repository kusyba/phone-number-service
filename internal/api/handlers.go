package api

import (
    "net/http"
    "strconv"
    "strings"
    
    "phone-number-service/internal/models"
    "phone-number-service/internal/service"
    
    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
)

type Handlers struct {
    phoneService *service.PhoneService
    groupService *service.GroupService
    validate     *validator.Validate
}

func NewHandlers(phoneService *service.PhoneService, groupService *service.GroupService) *Handlers {
    return &Handlers{
        phoneService: phoneService,
        groupService: groupService,
        validate:     validator.New(),
    }
}

func escapeHTML(s string) string {
    replacer := strings.NewReplacer(
        "&", "&amp;",
        "<", "&lt;",
        ">", "&gt;",
        "'", "&#39;",
        "\"", "&#34;",
    )
    return replacer.Replace(s)
}

func (h *Handlers) ImportNumbers(c echo.Context) error {
    var req models.ImportRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
    }
    
    req.Source = escapeHTML(req.Source)
    
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
    
    req.Number = escapeHTML(req.Number)
    req.Country = escapeHTML(req.Country)
    req.Region = escapeHTML(req.Region)
    req.Provider = escapeHTML(req.Provider)
    
    if err := h.validate.Struct(req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
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
    
    for i := range phones {
        phones[i].Number = escapeHTML(phones[i].Number)
        phones[i].Country = escapeHTML(phones[i].Country)
        phones[i].Region = escapeHTML(phones[i].Region)
        phones[i].Provider = escapeHTML(phones[i].Provider)
        phones[i].Source = escapeHTML(phones[i].Source)
    }
    
    return c.JSON(http.StatusOK, models.SearchResponse{
        Data:   phones,
        Total:  total,
        Limit:  filters.Limit,
        Offset: filters.Offset,
    })
}

func (h *Handlers) GetMe(c echo.Context) error {
    userIDHeader := c.Request().Header.Get("X-User-ID")
    if userIDHeader == "" {
        return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "user not authenticated"})
    }
    
    userID, err := strconv.Atoi(userIDHeader)
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid user id"})
    }
    
    groups, err := h.groupService.GetUserGroups(c.Request().Context(), userID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    
    flags := models.MergeFlags(groups)
    
    for i := range groups {
        groups[i].Name = escapeHTML(groups[i].Name)
        groups[i].Description = escapeHTML(groups[i].Description)
    }
    
    return c.JSON(http.StatusOK, models.MeResponse{
        UserID: userID,
        Groups: groups,
        Flags:  flags,
        Permissions: map[string]bool{
            "can_change_phone":             flags&models.FlagCanChangePhone != 0,
            "can_change_email":             flags&models.FlagCanChangeEmail != 0,
            "can_change_tariff":            flags&models.FlagCanChangeTariff != 0,
            "can_change_relative_settings": flags&models.FlagCanChangeRelativeSettings != 0,
            "can_leave_corporation":        flags&models.FlagCanLeaveCorporation != 0,
            "can_disable_incident_alerts":  flags&models.FlagCanDisableIncidentAlerts != 0,
            "can_delete_account":           flags&models.FlagCanDeleteAccount != 0,
        },
    })
}

func (h *Handlers) FormatPhoneByID(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid phone ID"})
    }
    
    formatted, countryCode, err := h.phoneService.FormatPhoneNumber(c.Request().Context(), id)
    if err != nil {
        return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Phone not found"})
    }
    
    return c.JSON(http.StatusOK, map[string]interface{}{
        "id":           id,
        "formatted":    escapeHTML(formatted),
        "country_code": escapeHTML(countryCode),
    })
}

func (h *Handlers) FormatPhoneByValue(c echo.Context) error {
    var req struct {
        Number string `json:"number" validate:"required"`
    }
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
    }
    
    req.Number = escapeHTML(req.Number)
    
    if err := h.validate.Struct(req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
    }
    
    formatted, countryCode, err := h.phoneService.FormatPhoneNumberByValue(c.Request().Context(), req.Number)
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
    }
    
    return c.JSON(http.StatusOK, map[string]interface{}{
        "original":     escapeHTML(req.Number),
        "formatted":    escapeHTML(formatted),
        "country_code": escapeHTML(countryCode),
    })
}
