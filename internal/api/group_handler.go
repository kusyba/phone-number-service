package api

import (
    "net/http"
    "strconv"
    
    "phone-number-service/internal/models"
    "phone-number-service/internal/service"
    
    "github.com/labstack/echo/v4"
)

type GroupHandlers struct {
    groupService *service.GroupService
}

func NewGroupHandlers(groupService *service.GroupService) *GroupHandlers {
    return &GroupHandlers{groupService: groupService}
}

func (h *GroupHandlers) CreateGroup(c echo.Context) error {
    var req struct {
        Name        string `json:"name" validate:"required"`
        Description string `json:"description"`
        Flags       int    `json:"flags"`
    }
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
    }
    if err := c.Validate(req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
    }
    
    group, err := h.groupService.CreateGroup(c.Request().Context(), req.Name, req.Description, req.Flags)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    return c.JSON(http.StatusCreated, group)
}

func (h *GroupHandlers) GetGroups(c echo.Context) error {
    var filters models.GroupFilters
    if err := c.Bind(&filters); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid filters"})
    }
    
    if err := c.Validate(filters); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
    }
    
    groups, total, err := h.groupService.GetGroups(c.Request().Context(), filters)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    
    return c.JSON(http.StatusOK, models.GroupListResponse{
        Data:   groups,
        Total:  total,
        Limit:  filters.Limit,
        Offset: filters.Offset,
    })
}

func (h *GroupHandlers) GetGroupByID(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid group id"})
    }
    
    group, err := h.groupService.GetGroupByID(c.Request().Context(), id)
    if err != nil {
        return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Group not found"})
    }
    return c.JSON(http.StatusOK, group)
}

func (h *GroupHandlers) UpdateGroup(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid group id"})
    }
    
    var req struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        Flags       int    `json:"flags"`
    }
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
    }
    
    if err := h.groupService.UpdateGroup(c.Request().Context(), id, req.Name, req.Description, req.Flags); err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func (h *GroupHandlers) DeleteGroup(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid group id"})
    }
    
    if err := h.groupService.DeleteGroup(c.Request().Context(), id); err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *GroupHandlers) AddUserToGroup(c echo.Context) error {
    groupID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid group id"})
    }
    userID, err := strconv.Atoi(c.Param("userId"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user id"})
    }
    
    if err := h.groupService.AddUserToGroup(c.Request().Context(), userID, groupID); err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    return c.JSON(http.StatusOK, map[string]string{"status": "user added to group"})
}

func (h *GroupHandlers) RemoveUserFromGroup(c echo.Context) error {
    groupID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid group id"})
    }
    userID, err := strconv.Atoi(c.Param("userId"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user id"})
    }
    
    if err := h.groupService.RemoveUserFromGroup(c.Request().Context(), userID, groupID); err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    return c.JSON(http.StatusOK, map[string]string{"status": "user removed from group"})
}

func (h *GroupHandlers) GetUserGroups(c echo.Context) error {
    userID, err := strconv.Atoi(c.Param("userId"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user id"})
    }
    
    groups, err := h.groupService.GetUserGroups(c.Request().Context(), userID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
    }
    return c.JSON(http.StatusOK, groups)
}
