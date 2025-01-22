package scheduler

import (
	"log"
	"time"

	"snappshop.ir/config"
	"snappshop.ir/internal/delivery"
)

type Scheduler struct {
	interval time.Duration
	service  *delivery.Service
}

func NewScheduler(cfg *config.Config, svc *delivery.Service) *Scheduler {
	if svc == nil {
		log.Fatalf("Service cannot be nil")
	}
	return &Scheduler{
		interval: time.Duration(cfg.SchedulerInterval) * time.Second,
		service:  svc,
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		s.run()
	}
}

func (s *Scheduler) run() {
	// Fetch scheduled requests and process them
	// Placeholder for scheduling logic
}
