package scheduler

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"snappshop.ir/internal/domain/entity"
	"snappshop.ir/internal/domain/repository"
	"snappshop.ir/internal/tpl"
)

type Dispatcher struct {
	repo      repository.OrderRepository
	tplClient tpl.Client
	// logger        logger.Logger
	logger        zerolog.Logger
	checkInterval time.Duration
	maxRetries    int
}

func NewDispatcher(
	repo repository.OrderRepository,
	tplClient tpl.Client,
	logger zerolog.Logger,
	checkInterval time.Duration,
	maxRetries int,
) *Dispatcher {
	return &Dispatcher{
		repo:          repo,
		tplClient:     tplClient,
		logger:        logger,
		checkInterval: checkInterval,
		maxRetries:    maxRetries,
	}
}

func (d *Dispatcher) Start(ctx context.Context) {
	d.logger.Info().
		Dur("check_interval", d.checkInterval).
		Msg("Starting order dispatcher")

	ticker := time.NewTicker(d.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.processOrders(ctx)
		case <-ctx.Done():
			d.logger.Info().Msg("Context canceled, stopping order dispatcher")
			return
		}
	}
}

func (d *Dispatcher) processOrders(ctx context.Context) {
	startTime := time.Now()
	log := d.logger.With().Str("method", "process_orders").Logger()

	log.Info().Msg("Checking for pending orders")

	orders, err := d.repo.GetByTimeToDeliver(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch orders")
		return
	}

	for _, order := range orders {
		orderLog := log.With().
			Str("order_number", order.OrderNumber).
			Time("deadline", order.TimeFrame.End).
			Logger()

		select {
		case <-ctx.Done():
			orderLog.Info().Msg("Context canceled, stopping order processing")
			return

		default:
			d.processOrder(ctx, &order, orderLog)
		}

	}

	log.Info().
		Dur("duration", time.Since(startTime)).
		Msg("Finished processing orders")
}

func (d *Dispatcher) processOrder(ctx context.Context, order *entity.Order, log zerolog.Logger) {
	log = log.With().Str("method", "processOrder").Logger()

	if time.Now().After(order.TimeFrame.End) {
		log.Warn().Msg("Order is already late")
		return
	}

	log.Info().
		Time("start_time", order.TimeFrame.Start).
		Msg("Processing order")

	for attempt := 1; attempt <= d.maxRetries; attempt++ {
		retryLog := log.With().Int("attempt", attempt).Logger()

		err := d.tplClient.CreateShipment(ctx, order.OrderNumber)
		if err == nil {
			if makeErr := d.repo.MarkOrderCompleted(ctx, order.OrderNumber); makeErr != nil {
				retryLog.Error().Err(makeErr).Msg("Failed to mark order as completed")
				return
			}

			retryLog.Info().Msg("Order successfully marked as completed")
			return
		}

		retryLog.Error().Err(err).Msg("Failed to process order")

		if attempt < d.maxRetries {
			retryDelay := time.Duration(attempt) * time.Second
			retryLog.Info().Dur("retry_delay", retryDelay).Msg("Scechduling retry")

			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				log.Info().Msg("Context canceled during retry delay")
				return
			}

		}
	}

	log.Error().
		Int("max_retries", d.maxRetries).
		Msg("Max retries reached, order processing failed")
}
