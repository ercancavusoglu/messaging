package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ercancavusoglu/messaging/internal/domain/message"
	"github.com/ercancavusoglu/messaging/internal/domain/scheduler"
	"github.com/ercancavusoglu/messaging/internal/infrastructure/cache"
	"github.com/ercancavusoglu/messaging/internal/infrastructure/persistance/postgres"
	"github.com/ercancavusoglu/messaging/internal/infrastructure/webhook"
	httpRouter "github.com/ercancavusoglu/messaging/internal/interfaces/http"
	"github.com/ercancavusoglu/messaging/internal/interfaces/http/handlers"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/messagingdb?sslmode=disable")
	if err != nil {
		fmt.Println("Error opening database", err)
		log.Fatal(err)
	}
	defer db.Close()

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	messageRepo := postgres.NewMessageRepository(db)
	webhookClient := webhook.NewClient("https://webhook.site/c3f13233-1ed4-429e-9649-8133b3b9c9cd", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo")
	cacheClient := cache.NewRedisAdapter(rdb)

	messageService := message.NewService(messageRepo, webhookClient, cacheClient)

	messageScheduler := scheduler.NewScheduler(messageService, 2*time.Second)

	// Start the scheduler in a separate goroutine
	ctx := context.Background()
	go func() {
		if err := messageScheduler.Start(ctx); err != nil {
			log.Printf("Scheduler error: %v", err)
		}
	}()

	messageHandler := handlers.NewMessageHandler(messageService, messageScheduler)

	router := httpRouter.NewRouter(messageHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
