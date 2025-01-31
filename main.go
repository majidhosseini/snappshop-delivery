package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"snappshop.ir/cmd"
	"snappshop.ir/config"

	// "snappshop.ir/internal/delivery"
	"snappshop.ir/internal/domain/repository"
	"snappshop.ir/internal/scheduler"
	"snappshop.ir/internal/tpl"
	"snappshop.ir/pkg/http"
	"snappshop.ir/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := sqlx.Connect("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	cmd.RunMigrations(db, "./db/migrations")

	mainDB, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "order_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Initialize dependencies
	orderRepository := repository.NewOrderRepository(mainDB)

	tplClient := tpl.NewClient("test.ir", "test_token")
	logger := logger.New("delivery-service")

	// deliveryRepository := delivery.NewGORMRepository(mainDB)
	// deliveryService := delivery.NewService(orderRepository)

	schedulerService := scheduler.NewDispatcher(orderRepository, tplClient, logger.Logger, cfg.SchedulerInterval, 3)
	go schedulerService.Start(ctx)

	// Start HTTP server
	server := http.NewServer(cfg.HTTP.Port)

	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(shutdownCtx)
}
