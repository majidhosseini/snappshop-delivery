package delivery

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"snappshop.ir/internal/domain/entity"
	"snappshop.ir/internal/domain/repository"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerService struct {
	orderRepository    repository.OrderRepository
	deliveryRepository repository.DeliveryRepository
	kafkaConsumer      *kafka.Consumer
}

func NewConsumerService(orderRepo repository.OrderRepository, deliveryRepo repository.DeliveryRepository, kafkaConsumer *kafka.Consumer) *ConsumerService {
	return &ConsumerService{
		orderRepository:    orderRepo,
		deliveryRepository: deliveryRepo,
		kafkaConsumer:      kafkaConsumer,
	}
}

func (s *ConsumerService) ConsumeMessages(ctx context.Context, topic string) {
	err := s.kafkaConsumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("failed to subscribe to topic: %v", err)
	}

	for {
		msg, err := s.kafkaConsumer.ReadMessage(-1)
		if err != nil {
			log.Printf("consumer error: %v (%v)\n", err, msg)
			continue
		}

		var order entity.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		if err := s.handleOrder(ctx, &order); err != nil {
			log.Printf("failed to handle order: %v", err)
		}
	}
}

func (s *ConsumerService) handleOrder(ctx context.Context, order *entity.Order) error {
	// Create a new delivery
	delivery := &entity.Delivery{
		OrderID:        order.ID,
		Provider:       "3PL Provider", // Example provider name
		Origin:         order.Origin,
		Destination:    order.Destination,
		TimeFrameStart: order.TimeFrame.Start,
		TimeFrameEnd:   order.TimeFrame.End,
		Status:         entity.StateInit,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.deliveryRepository.Create(delivery); err != nil {
		return err
	}

	// Update the order status
	order.Status = entity.StatusInProgress
	order.UpdatedAt = time.Now()
	return s.orderRepository.Update(order)
}
