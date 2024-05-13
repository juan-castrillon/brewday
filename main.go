package main

import (
	"brewday/internal/app"
	"brewday/internal/config"
	"brewday/internal/notifications"
	"brewday/internal/render"
	"brewday/internal/store/memory"
	recipe_store_sql "brewday/internal/store/sql"
	summary_store_memory "brewday/internal/summary/memory"
	tl_store_memory "brewday/internal/timeline/memory"
	tl_store_sql "brewday/internal/timeline/sql"
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

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
	switch config.Store.StoreType {
	case "sql":
		db, err := sql.Open("sqlite3", "file:"+config.Store.Path+"?_foreign_keys=true")
		if err != nil {
			log.Fatal().Err(err).Msg("Error while initializing db store")
		}
		defer db.Close()
		s, err := recipe_store_sql.NewPersistentStore(db)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while initializing db store")
		}
		defer s.Close()
		components.Store = s
		tls, err := tl_store_sql.NewTimelinePersistentStore(db)
		if err != nil {
			log.Fatal().Err(err).Msg("Error while initializing db store")
		}
		defer tls.Close()
		components.TL = tls
	case "memory":
		components.Store = memory.NewMemoryStore()
		components.TL = tl_store_memory.NewTimelineMemoryStore()
	default:
		log.Fatal().Msg("Invalid store type")
	}
	components.SummaryStore = summary_store_memory.NewSummaryMemoryStore()
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
