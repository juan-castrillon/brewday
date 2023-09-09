package main

import (
	"brewday/internal/app"
	"context"
	"embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

//go:embed web
var staticFS embed.FS

var version = "" // This is set by the build process

func main() {
	app, err := app.NewApp(staticFS)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while initializing the app")
	}
	log.Info().Msgf("Starting BrewDay version %s", version)
	go func() {
		if err := app.Run(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Error while running the app")
		}
	}()
	// Graceful shutdown with 10 seconds timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Gracefully shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Fatal().Err(err).Msg("Error while shutting down the app")
	}
	log.Info().Msg("Server shutdown complete")
	os.Exit(0)
}
