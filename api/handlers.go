package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pindamonhangaba/apiculi"
	"github.com/pindamonhangaba/urlshorts/service"
)

// handleRedirect handles the redirection of shortened URLs
func (s *Server) handleRedirect(c *apiculi.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, service.ErrorResponse{
			Error: "Missing URL code",
		})
	}

	// Get the URL from the database
	url, err := s.config.DB.GetURL(code)
	if err != nil {
		return c.JSON(http.StatusNotFound, service.ErrorResponse{
			Error: "URL not found",
		})
	}

	// Update visit count
	url.Visits++
	// We're ignoring the error here for simplicity, but in a production app
	// you'd want to handle this error properly
	s.config.DB.SaveURL(url)

	// Redirect to the original URL
	return c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// handleCreateURL handles the creation of a new shortened URL
func (s *Server) handleCreateURL(c *apiculi.Context) error {
	// Validate API key
	apiKey := c.Request().Header.Get("X-API-Key")
	if !service.ValidateAPIKey(apiKey, s.config.APIKey) {
		return c.JSON(http.StatusUnauthorized, service.ErrorResponse{
			Error: "Invalid API key",
		})
	}

	// Parse request body
	var req service.CreateURLRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, service.ErrorResponse{
			Error: "Invalid request payload",
		})
	}

	// Validate the original URL
	if req.OriginalURL == "" {
		return c.JSON(http.StatusBadRequest, service.ErrorResponse{
			Error: "Original URL is required",
		})
	}

	// Clean pretty name (if provided)
	prettyName := strings.TrimSpace(req.PrettyName)

	// Generate a random code
	code, err := service.GenerateRandomCode(service.DefaultCodeLength)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, service.ErrorResponse{
			Error: "Failed to generate code",
		})
	}

	// Create the URL model
	url := &service.URL{
		Code:        code,
		OriginalURL: req.OriginalURL,
		PrettyName:  prettyName,
		CreatedAt:   time.Now(),
		Visits:      0,
	}

	// Save the URL to the database
	if err := s.config.DB.SaveURL(url); err != nil {
		return c.JSON(http.StatusInternalServerError, service.ErrorResponse{
			Error: "Failed to save URL",
		})
	}

	// Build the short URL
	shortURL := service.BuildShortURL(s.config.BaseURL, code, prettyName)

	// Return the response
	return c.JSON(http.StatusCreated, service.CreateURLResponse{
		Code:        code,
		ShortURL:    shortURL,
		OriginalURL: req.OriginalURL,
		PrettyName:  prettyName,
	})
}

// handleListURLs handles listing all URLs
func (s *Server) handleListURLs(c *apiculi.Context) error {
	// Validate API key
	apiKey := c.Request().Header.Get("X-API-Key")
	if !service.ValidateAPIKey(apiKey, s.config.APIKey) {
		return c.JSON(http.StatusUnauthorized, service.ErrorResponse{
			Error: "Invalid API key",
		})
	}

	// Get all URLs from the database
	urls, err := s.config.DB.ListURLs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, service.ErrorResponse{
			Error: "Failed to retrieve URLs",
		})
	}

	return c.JSON(http.StatusOK, urls)
}
