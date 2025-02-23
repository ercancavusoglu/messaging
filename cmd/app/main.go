package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ercancavusoglu/messaging/internal/adapters"
)

func main() {
	container, err := adapters.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	logger := container.Logger
	logger.Info("=== Starting Messaging Application ===")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := container.Start(ctx); err != nil {
		logger.Fatalf("Failed to start services: %v", err)
	}

	logger.Info("[Main] Application started successfully!")
	logger.Info("Press Ctrl+C to shutdown...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("\n[Main] Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := container.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("[Main] Shutdown error: %v", err)
	}

	logger.Info("[Main] Shutdown complete")
}
