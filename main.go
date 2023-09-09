package main

import (
	"brewday/internal/app"
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed web
var staticFS embed.FS

func main() {
	app, err := app.NewApp(staticFS)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := app.Run(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	// Graceful shutdown with 10 seconds timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Gracefully shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Fatalf("Error while shutting down the server. Error %s", err.Error())
	}
	log.Printf("Server shutdown complete")
	os.Exit(0)
}
