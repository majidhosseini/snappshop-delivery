package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnectionString string
	HTTPPort           string
	SchedulerInterval  int
	ServerAddress      string
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
		DBConnectionString: getEnv("DB_CONNECTION_STRING"),
		HTTPPort:           getEnv("HTTP_PORT"),
		SchedulerInterval:  schedulerInterval,
		ServerAddress:      getEnv("SERVER_ADDRESS"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
