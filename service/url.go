package service

import (
	"time"
)

// URL represents a shortened URL
type URL struct {
	Code       string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	PrettyName string    `json:"pretty_name,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	Visits     int       `json:"visits"`
}

// CreateURLRequest represents a request to create a shortened URL
type CreateURLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	PrettyName  string `json:"pretty_name,omitempty"`
}

// CreateURLResponse represents the response after creating a shortened URL
type CreateURLResponse struct {
	Code       string `json:"code"`
	ShortURL   string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	PrettyName string `json:"pretty_name,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}