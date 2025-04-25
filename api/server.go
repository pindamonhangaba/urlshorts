package api

import (
	"fmt"
	"net/http"

	"github.com/pindamonhangaba/apiculi"
	"github.com/pindamonhangaba/urlshorts/db"
)

// ServerConfig contains configuration for the API server
type ServerConfig struct {
	Host    string
	Port    string
	DB      *db.DB
	APIKey  string
	BaseURL string
}

// Server represents the API server
type Server struct {
	config ServerConfig
	api    *apiculi.API
}

// NewServer creates a new API server
func NewServer(config ServerConfig) *Server {
	return &Server{
		config: config,
		api:    apiculi.New(),
	}
}

// Start starts the API server
func (s *Server) Start() error {
	// Register routes
	s.registerRoutes()

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	// Create the HTTP handler
	http.HandleFunc("/", s.handleHTTP)

	// Start the server
	return http.ListenAndServe(addr, nil)
}

// registerRoutes registers all the API routes
func (s *Server) registerRoutes() {
	// Register the handlers with apiculi
	s.api.GET("/{code}", s.handleRedirect)
	s.api.GET("/{code}/{prettyname}", s.handleRedirect)
	s.api.POST("/api/urls", s.handleCreateURL)
	s.api.GET("/api/urls", s.handleListURLs)
}

// handleHTTP handles all HTTP requests and routes them to the appropriate handler
func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	s.api.ServeHTTP(w, r)
}
