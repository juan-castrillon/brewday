package main

import (
	"brewday/internal/app"
	"brewday/internal/config"
	"brewday/internal/notifications"
	"brewday/internal/render"
	"brewday/internal/store/memory"
	"brewday/internal/store/sql"
	summaryrecorder "brewday/internal/summary_recorder"
	"brewday/internal/timeline"
	"context"
	"embed"
	"flag"
	"fmt"
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
	// Parse configuration
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the config file")
	flag.Parse()
	config, err := config.LoadConfig(configPath) // If config path is empty, it will attempt to load from the environment
	if err != nil {
		log.Fatal().Err(err).Msg("Error while loading config")
	}
	// Configure application
	runningPort := fmt.Sprintf(":%d", config.App.Port)
	components := &app.AppComponents{}
	// Initialize components
	components.Renderer = render.NewTemplateRenderer()
	components.TL = timeline.NewTimelineStore()
	var store app.RecipeStore
	switch config.Store.StoreType {
	case "sql":
		s, err := sql.NewPersistentStore(config.Store.Path)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while initializing db store")
		}
		defer s.Close()
		store = s
	case "memory":
		store = memory.NewMemoryStore()
	default:
		log.Fatal().Msg("Invalid store type")
	}
	components.Store = store
	components.SummaryStore = summaryrecorder.NewSummaryRecorderStore()
	if config.Notification.Enabled {
		n, err := notifications.NewGotifyNotifier(
			config.Notification.GotifyURL,
			config.Notification.Username,
			config.Notification.Password,
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while initializing notifier")
		}
		components.Notifier = n
	}
	app, err := app.NewApp(staticFS, components)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while initializing the app")
	}
	log.Info().Msgf("Starting BrewDay version %s", version)
	go func() {
		if err := app.Run(runningPort); err != nil && err != http.ErrServerClosed {
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
