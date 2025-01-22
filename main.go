package main

import (
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"snappshop.ir/config"
	"snappshop.ir/internal/delivery"
	"snappshop.ir/internal/scheduler"
	"snappshop.ir/migrations"

	"snappshop.ir/seeders"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := sqlx.Connect("postgres", cfg.DBConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	migrations.RunMigrations(db, "./migrations")

	mainDB, err := gorm.Open(postgres.Open(cfg.DBConnectionString), &gorm.Config{})

	// Initialize dependencies
	deliveryRepository := delivery.NewGORMRepository(mainDB)
	deliveryService := delivery.NewService(deliveryRepository)
	schedulerService := scheduler.NewScheduler(cfg, deliveryService)

	// Seed initial data
	seeders.SeedRequests(deliveryService, 100)
	go seeders.SeedScheduledRequests(deliveryService, 30*time.Second)

	// Start scheduler
	go schedulerService.Start()

	// Start HTTP server
	if err := StartHTTPServer(cfg, deliveryService); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// StartHTTPServer starts the HTTP server
func StartHTTPServer(cfg *config.Config, deliveryService *delivery.Service) error {

	http.HandleFunc("/deliveries", func(w http.ResponseWriter, r *http.Request) {
		// Handle delivery requests
		w.Write([]byte("Delivery service is running"))
	})

	log.Printf("Starting HTTP server on %s", cfg.ServerAddress)

	return http.ListenAndServe(cfg.ServerAddress, nil)
}
