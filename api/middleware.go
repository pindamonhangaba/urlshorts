package api

import (
	"log"
	"net/http"
	"time"

	"github.com/pindamonhangaba/apiculi"
	"github.com/pindamonhangaba/urlshorts/service"
)

// AuthMiddleware provides API key authentication for protected routes
func (s *Server) AuthMiddleware(next apiculi.Handler) apiculi.Handler {
	return func(c *apiculi.Context) error {
		apiKey := c.Request().Header.Get("X-API-Key")
		if !service.ValidateAPIKey(apiKey, s.config.APIKey) {
			return c.JSON(http.StatusUnauthorized, service.ErrorResponse{
				Error: "Invalid API key",
			})
		}
		return next(c)
	}
}

// LoggingMiddleware logs information about incoming requests
func LoggingMiddleware(next apiculi.Handler) apiculi.Handler {
	return func(c *apiculi.Context) error {
		start := time.Now()

		err := next(c)

		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s - %v",
			c.Request().Method,
			c.Request().URL.Path,
			duration,
			err,
		)

		return err
	}
}
