package main

import (
	"brewday/internal/app"
	"brewday/internal/config"
	"brewday/internal/notifications"
	"brewday/internal/render"
	"brewday/internal/store/memory"
	"brewday/internal/summary_recorder/markdown"
	"brewday/internal/timeline/basic"
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

const (
	defaultConfigPath = "/etc/brewday/config.yaml"
)

func main() {
	// Parse configuration
	var configPath string
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to the config file")
	flag.Parse()
	if configPath == defaultConfigPath {
		log.Warn().Msgf("Using default config path %s", defaultConfigPath)
	} else {
		log.Info().Msgf("Using config path %s", configPath)
	}
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while loading config")
	}
	// Configure application
	runningPort := fmt.Sprintf(":%d", config.App.Port)
	components := &app.AppComponents{}
	// Initialize components
	components.Renderer = render.NewTemplateRenderer()
	components.TL = basic.NewBasicTimeline()
	components.Store = memory.NewMemoryStore()
	components.Summary = markdown.NewMarkdownSummaryRecorder()
	if config.Notification.Enabled {
		components.Notifier = notifications.NewGotifyNotifier(config.Notification.AppToken, config.Notification.GotifyURL)
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
