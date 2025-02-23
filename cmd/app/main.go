package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ercancavusoglu/messaging/internal/adapters"
	"github.com/ercancavusoglu/messaging/internal/adapters/scheduler"

	"github.com/ercancavusoglu/messaging/internal/adapters/eventbus"
	"github.com/ercancavusoglu/messaging/internal/adapters/persistance/cache"
	"github.com/ercancavusoglu/messaging/internal/adapters/persistance/postgres"
	"github.com/ercancavusoglu/messaging/internal/adapters/webhook"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("=== Starting Messaging Application ===")

	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/messagingdb?sslmode=disable")
	if err != nil {
		fmt.Println("Error opening database", err)
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("[Main] PostgreSQL connection established")

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()
	fmt.Println("[Main] Redis connection established")

	eventBus, err := eventbus.NewRabbitMQEventBus("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer eventBus.Close()
	fmt.Println("[Main] RabbitMQ connection established")

	fmt.Println("[Main] Initializing services...")
	messageRepo := postgres.NewMessageRepository(db)
	webhookClient := webhook.NewClient("https://webhook.site/ae17c131-349d-410b-8cc5-2f17c823ccca", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo")
	cacheClient := cache.NewRedisAdapter(rdb)

	messageService := adapters.NewMessageService(messageRepo, webhookClient, cacheClient, eventBus)
	messageScheduler := scheduler.NewSchedulerService(messageService, 2*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := messageScheduler.Start(ctx); err != nil {
			log.Printf("[Main] Scheduler error: %v", err)
		}
	}()

	fmt.Println("[Main] Starting HTTP server...")
	messageHandler := adapters.NewMessageHandler(messageService, messageScheduler)
	router := adapters.NewRouter(messageHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Println("[Main] HTTP server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Main] Server error: %v", err)
		}
	}()

	fmt.Println("[Main] Application started successfully!")
	fmt.Println("Press Ctrl+C to shutdown...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n[Main] Shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[Main] Server shutdown error: %v", err)
	}
	fmt.Println("[Main] Shutdown complete")
}
