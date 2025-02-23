package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ercancavusoglu/messaging/internal/adapters/consumer"
	"github.com/ercancavusoglu/messaging/internal/adapters/eventbus"
	logrus "github.com/ercancavusoglu/messaging/internal/adapters/logger"
	"github.com/ercancavusoglu/messaging/internal/adapters/persistance/cache"
	"github.com/ercancavusoglu/messaging/internal/adapters/persistance/postgres"
	"github.com/ercancavusoglu/messaging/internal/adapters/scheduler"
	"github.com/ercancavusoglu/messaging/internal/adapters/webhook"
	"github.com/ercancavusoglu/messaging/internal/ports"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Container struct {
	DB             *sql.DB
	Redis          *redis.Client
	EventBus       *eventbus.RabbitMQEventBus
	MessageService ports.MessageService
	Scheduler      *scheduler.SchedulerService
	Consumer       *consumer.Consumer
	Server         *http.Server
	Logger         ports.Logger
}

func NewContainer() (*Container, error) {
	// Initialize logger first
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "consumer.log"
	}
	logger, err := logrus.NewLogrusAdapter(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Warnf("Warning: .env file not found or error loading it: %v", err)
	}

	// Initialize database
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	logger.Info("Connecting to database...")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	logger.Info("Database connection established")

	// Initialize Redis
	logger.Info("Connecting to Redis...")
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
	})
	logger.Info("Redis connection established")

	// Initialize RabbitMQ
	logger.Info("Connecting to RabbitMQ...")
	eventBus, err := eventbus.NewRabbitMQEventBus(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	logger.Info("RabbitMQ connection established")

	// Initialize repositories and services
	logger.Info("Initializing services...")
	messageRepo := postgres.NewMessageRepository(db)
	webhookClient := webhook.NewClient(os.Getenv("WEBHOOK_URL"), os.Getenv("WEBHOOK_TOKEN"))
	cacheClient := cache.NewRedisAdapter(rdb)
	messageSvc := NewMessageService(messageRepo, webhookClient, cacheClient, eventBus)
	messageScheduler := scheduler.NewSchedulerService(messageSvc, 2*time.Second, logger)
	messageConsumer := consumer.NewConsumer(webhookClient, messageRepo, cacheClient, eventBus, 5, logger)

	// Initialize HTTP server
	messageHandler := NewMessageHandler(messageSvc, messageScheduler)
	router := NewRouter(messageHandler)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: router,
	}

	return &Container{
		DB:             db,
		Redis:          rdb,
		EventBus:       eventBus,
		MessageService: messageSvc,
		Scheduler:      messageScheduler,
		Consumer:       messageConsumer,
		Server:         server,
		Logger:         logger,
	}, nil
}

func (c *Container) Start(ctx context.Context) error {
	// Start consumer
	if err := c.Consumer.Start(); err != nil {
		c.Logger.Errorf("[Container] Consumer error: %v", err)
		return err
	}

	// Start scheduler
	go func() {
		if err := c.Scheduler.Start(ctx); err != nil {
			c.Logger.Errorf("[Container] Scheduler error: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		c.Logger.Infof("[Container] HTTP server listening on %s", c.Server.Addr)
		if err := c.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.Logger.Errorf("[Container] Server error: %v", err)
		}
	}()

	return nil
}

func (c *Container) Shutdown(ctx context.Context) error {
	c.Logger.Info("[Container] Shutting down...")

	// Stop consumer and scheduler
	c.Consumer.Stop()
	c.Scheduler.Stop()

	// Shutdown HTTP server
	if err := c.Server.Shutdown(ctx); err != nil {
		c.Logger.Errorf("Failed to shutdown server: %v", err)
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	// Close connections
	if err := c.DB.Close(); err != nil {
		c.Logger.Errorf("Failed to close database connection: %v", err)
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	if err := c.Redis.Close(); err != nil {
		c.Logger.Errorf("Failed to close Redis connection: %v", err)
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	if err := c.EventBus.Close(); err != nil {
		c.Logger.Errorf("Failed to close event bus connection: %v", err)
		return fmt.Errorf("failed to close event bus connection: %w", err)
	}

	c.Logger.Info("[Container] Shutdown complete")
	return nil
}
