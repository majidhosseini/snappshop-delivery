package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	HTTP struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
	DB struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
	}
	SchedulerInterval time.Duration
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	schedulerInterval, err := strconv.Atoi(os.Getenv("SCHEDULER_INTERVAL"))
	if err != nil {
		log.Printf("Invalid SCHEDULER_INTERVAL, using default 10 seconds: %v", err)
		schedulerInterval = 10
	}

	return &Config{
		Env: "dev",
		HTTP: struct {
			Port         string
			ReadTimeout  time.Duration
			WriteTimeout time.Duration
		}{
			Port:         "8080",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		DB: struct {
			DSN          string
			MaxOpenConns int
			MaxIdleConns int
		}{
			DSN:          "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
			MaxOpenConns: 25,
			MaxIdleConns: 25,
		},
		SchedulerInterval: time.Duration(schedulerInterval) * time.Second,
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
