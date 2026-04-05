package models

import "time"

type Phone struct {
    ID        int       `json:"id"`
    Number    string    `json:"number"`
    Country   string    `json:"country"`
    Region    string    `json:"region"`
    Provider  string    `json:"provider"`
    Source    string    `json:"source"`
    CreatedAt time.Time `json:"created_at"`
}

type ImportRequest struct {
    Numbers []string `json:"numbers" validate:"required,min=1"`
    Source  string   `json:"source" validate:"required"`
}

type ImportResponse struct {
    Accepted int `json:"accepted"`
    Skipped  int `json:"skipped"`
    Errors   int `json:"errors"`
}

type SearchRequest struct {
    Number   string `query:"number"`
    Country  string `query:"country"`
    Region   string `query:"region"`
    Provider string `query:"provider"`
    Limit    int    `query:"limit"`
    Offset   int    `query:"offset"`
}

type SearchResponse struct {
    Data   []Phone `json:"data"`
    Total  int64   `json:"total"`
    Limit  int     `json:"limit"`
    Offset int     `json:"offset"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
