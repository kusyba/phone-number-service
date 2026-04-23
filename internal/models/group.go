package models

import "time"

const (
    FlagCanChangePhone            = 1 << 4
    FlagCanChangeEmail            = 1 << 5
    FlagCanChangeTariff           = 1 << 6
    FlagCanChangeRelativeSettings = 1 << 7
    FlagCanLeaveCorporation       = 1 << 8
    FlagCanDisableIncidentAlerts  = 1 << 9
    FlagCanDeleteAccount          = 1 << 10
)

type Group struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Flags       int       `json:"flags"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type GroupFilters struct {
    Name   string `query:"name"`
    Search string `query:"search"`
    Sort   string `query:"sort"`
    Limit  int    `query:"limit" validate:"min=1,max=100"`
    Offset int    `query:"offset" validate:"min=0"`
}

type GroupListResponse struct {
    Data   []Group `json:"data"`
    Total  int64   `json:"total"`
    Limit  int     `json:"limit"`
    Offset int     `json:"offset"`
}

type MeResponse struct {
    UserID      int            `json:"user_id"`
    Groups      []Group        `json:"groups"`
    Flags       int            `json:"flags"`
    Permissions map[string]bool `json:"permissions"`
}

func MergeFlags(groups []Group) int {
    flags := 0
    for _, g := range groups {
        flags |= g.Flags
    }
    return flags
}
