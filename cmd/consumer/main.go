package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ercancavusoglu/messaging/internal/adapters"
)

func main() {
	container, err := adapters.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	container.Logger.Info("=== Starting Message Consumer ===")

	if err := container.Consumer.Start(); err != nil {
		container.Logger.Fatalf("[Consumer] Failed to start consumer: %v", err)
	}

	container.Logger.Info("[Consumer] Started successfully!")
	container.Logger.Info("Press Ctrl+C to shutdown...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	container.Logger.Info("\n[Consumer] Shutting down...")
	container.Consumer.Stop()
	container.Logger.Info("[Consumer] Shutdown complete")
}
