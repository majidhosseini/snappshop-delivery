package seeders

import (
	"delivery-service/internal/delivery"
	"testing"
)

func TestSeedRequests(t *testing.T) {
	repo := delivery.NewMockRepository()
	service := delivery.NewService(repo)

	SeedRequests(service, 10)

	if len(repo.requests) != 10 {
		t.Fatalf("Expected 10 requests, got %d", len(repo.requests))
	}
}

func TestSeedScheduledRequests(t *testing.T) {
	repo := delivery.NewMockRepository()
	service := delivery.NewService(repo)

	go SeedScheduledRequests(service, 1) // Run for a short time
}
