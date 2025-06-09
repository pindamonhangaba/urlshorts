package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pindamonhangaba/apiculi/endpoint"
	"github.com/pindamonhangaba/urlshorts/db"
	"github.com/pindamonhangaba/urlshorts/service"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ServerConfig contains configuration for the API server
type ServerConfig struct {
	DB      *db.DB
	APIKey  string
	BaseURL string
}

// HTTPServer represents the API server
type HTTPServer struct {
	config ServerConfig
	stub   bool
}

// NewServer creates a new API server
func NewServer(config ServerConfig) *HTTPServer {
	return &HTTPServer{
		config: config,
	}
}

func (h *HTTPServer) validateDeps() error {
	if h.config.DB == nil || h.config.APIKey == "" {
		return errors.New("missing dependency")
	}
	return nil
}

// Run create a new echo server
func (h *HTTPServer) Register(e *echo.Echo) (*endpoint.OpenAPI, error) {
	err := h.validateDeps()
	if !h.stub && err != nil {
		return nil, err
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, errors.Wrap(err, "logger init")
	}

	e.Use(ZapLogger(zapLogger))

	gAuth := e.Group("/admin")
	gAuth.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:Authorization",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == h.config.APIKey, nil
		},
	}))

	oapi := endpoint.NewOpenAPI("URL Shorts", "v1")
	oapi.AddServer(e.Server.Addr, h.config.BaseURL)

	e.Add(endpoint.EchoWithContext(
		endpoint.Get("/:slug/:code"),
		oapi.Route("url.Redirect", `Redirect to the original URL`),
		func(in endpoint.EndpointInput[any, struct {
			Slug string `json:"slug"`
			Code string `json:"code"`
		}, struct {
			Context string `json:"context,omitempty"`
		}, any], c echo.Context) (res endpoint.DataResponse[endpoint.SingleItemData[any]], err error) {

			// Get the URL from the database
			url, err := h.config.DB.GetURL(in.Params.Code)
			if err != nil {
				return res, errors.New("URL not found")
			}

			// Update visit count
			url.Visits++
			// We're ignoring the error here for simplicity, but in a production app
			// you'd want to handle this error properly
			err = h.config.DB.SaveURL(url)
			if err != nil {
				return res, err
			}
			// Redirect to the original URL
			c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
			return res, nil
		},
	))

	gAuth.Add(endpoint.Echo(
		endpoint.Post("/shorts"),
		oapi.Route("url.Shorts", `Create a short URL for a given large URL`),
		func(in endpoint.EndpointInput[any, any, struct {
			Context string `json:"context,omitempty"`
		}, struct {
			URL        string  `json:"url"`
			PrettyName *string `json:"pretty_name,omitempty"`
		}]) (res endpoint.DataResponse[endpoint.SingleItemData[any]], err error) {

			prettyName := ""
			if in.Body.PrettyName == nil {
				prettyName = strings.TrimSpace(*in.Body.PrettyName)
			}

			code, err := service.GenerateRandomCode(service.DefaultCodeLength)
			if err != nil {
				return res, err
			}

			url := &service.URL{
				Code:        code,
				OriginalURL: in.Body.URL,
				PrettyName:  prettyName,
				CreatedAt:   time.Now(),
				Visits:      0,
			}

			if err := h.config.DB.SaveURL(url); err != nil {
				return res, err
			}

			shortURL := service.BuildShortURL(h.config.BaseURL, code, prettyName)

			res.Context = in.Query.Context
			res.Data.Item = service.CreateURLResponse{
				Code:        code,
				ShortURL:    shortURL,
				OriginalURL: in.Body.URL,
				PrettyName:  prettyName,
			}
			return res, nil
		},
	))

	swagapijson, err := oapi.T().MarshalJSON()
	if err != nil {
		panic(err)
	}
	e.GET("/docs/swagger.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, json.RawMessage(swagapijson))
	})
	e.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Redoc</title>
				<!-- needed for adaptive design -->
				<meta charset="utf-8"/>
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
				<!--
				Redoc doesn't change outer page styles
				-->
				<style>
				body {
					margin: 0;
					padding: 0;
				}
				</style>
			</head>
			<body>
				<redoc spec-url='/docs/swagger.json'></redoc>

				<script src="https://cdn.jsdelivr.net/npm/redoc@latest/bundles/redoc.standalone.js"> </script>
			</body>
			</html>
		`)
	})
	return &oapi, nil
}
