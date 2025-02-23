package main

import (
	"database/sql"
	"fmt"
	"github.com/ercancavusoglu/messaging/internal/adapters/consumer"
	"github.com/ercancavusoglu/messaging/internal/adapters/persistance/cache"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ercancavusoglu/messaging/internal/adapters/eventbus"
	"github.com/ercancavusoglu/messaging/internal/adapters/persistance/postgres"
	"github.com/ercancavusoglu/messaging/internal/adapters/webhook"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("=== Starting Message Consumer ===")

	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/messagingdb?sslmode=disable")
	if err != nil {
		fmt.Println("Error opening database", err)
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("[Consumer] PostgreSQL connection established")

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()
	fmt.Println("[Consumer] Redis connection established")

	eventBus, err := eventbus.NewRabbitMQEventBus("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer eventBus.Close()
	fmt.Println("[Consumer] RabbitMQ connection established")

	fmt.Println("[Consumer] Initializing services...")
	messageRepo := postgres.NewMessageRepository(db)
	webhookClient := webhook.NewClient("https://webhook.site/ae17c131-349d-410b-8cc5-2f17c823ccca", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo")
	cacheClient := cache.NewRedisAdapter(rdb)

	consumerClient := consumer.NewConsumer(webhookClient, messageRepo, cacheClient, eventBus, 5)
	if err := consumerClient.Start(); err != nil {
		log.Fatalf("[Consumer] Failed to start consumer: %v", err)
	}

	fmt.Println("[Consumer] Started successfully!")
	fmt.Println("Press Ctrl+C to shutdown...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n[Consumer] Shutting down...")
	consumerClient.Stop()
	fmt.Println("[Consumer] Shutdown complete")
}
