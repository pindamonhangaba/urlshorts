package main

import (
	"log"

	"github.com/pindamonhangaba/commandments"
	"github.com/pindamonhangaba/urlshorts/api"
	"github.com/pindamonhangaba/urlshorts/db"
)

type envArgs struct {
	Host    string `env:"HOST" default:"localhost"`
	Port    string `env:"PORT" default:"8080"`
	DBPath  string `env:"DB_PATH" default:"shortener.db"`
	APIKey  string `env:"API_KEY" default:"your-api-key-here"`
	BaseURL string `env:"BASE_URL" default:"http://localhost:8080"`
}

func main() {
	_ = commandments.MustCMD("link-shortener", commandments.WithConfig(
		func(config envArgs) error {
			// Initialize the database
			database, err := db.NewDB(config.DBPath)
			if err != nil {
				return err
			}
			defer database.Close()

			// Initialize and start the API server
			server := api.NewServer(api.ServerConfig{
				Host:    config.Host,
				Port:    config.Port,
				DB:      database,
				APIKey:  config.APIKey,
				BaseURL: config.BaseURL,
			})

			log.Printf("Server starting on %s:%s", config.Host, config.Port)
			return server.Start()
		}), commandments.WithDefaultConfig(envArgs{
		Host:    "localhost",
		Port:    "8080",
		DBPath:  "shortener.db",
		APIKey:  "your-api-key-here",
		BaseURL: "http://localhost:8080",
	}))
}
