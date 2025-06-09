package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pindamonhangaba/commandments"
	"github.com/pindamonhangaba/urlshorts/api"
	"github.com/pindamonhangaba/urlshorts/db"
	"github.com/spf13/cobra"
)

type envArgs struct {
	Host    string `flag:"host" default:"localhost"`
	Port    string `flag:"port" default:"8080"`
	DBPath  string `flag:"db_path" default:"shortener.db"`
	APIKey  string `flag:"api_key" default:"your-api-key-here"`
	BaseURL string `flag:"base_url" default:"http://localhost:8080"`
}

func main() {
	serveCmd := commandments.MustCMD("serve", commandments.WithConfig(
		func(config envArgs) error {
			// Initialize the database
			database, err := db.NewDB(config.DBPath)
			if err != nil {
				return err
			}
			defer database.Close()

			e := echo.New()

			// Initialize and start the API server
			serv := api.NewServer(api.ServerConfig{
				DB:      database,
				APIKey:  config.APIKey,
				BaseURL: config.BaseURL,
			})

			serv.Register(e)

			start := func() error { return e.Start(config.Host + ":" + config.Port) }

			go func() {
				err := start()
				if err != nil && err != http.ErrServerClosed {
					e.Logger.Fatalf("server shutdown %s", err)
				} else {
					e.Logger.Fatal("shutting down the server")
				}
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return e.Shutdown(ctx)
		}), commandments.WithDefaultConfig(envArgs{
		Host:    "localhost",
		Port:    "8080",
		DBPath:  "shortener.db",
		APIKey:  "your-api-key-here",
		BaseURL: "https://relatorio.link",
	}))

	cmd := &cobra.Command{
		Use: "urlshorts",
	}
	cmd.PersistentFlags().String("config", "", "Path for config file")
	cmd.AddCommand(serveCmd)
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
