package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"snappshop.ir/internal/domain/entity"
	"snappshop.ir/internal/domain/repository"
)

type Service struct {
	repo          repository.OrderRepository
	kafkaProducer *kafka.Producer
	topic         string
}

func NewService(repo repository.OrderRepository, kafkaProducer *kafka.Producer, topic string) *Service {
	return &Service{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		topic:         topic,
	}
}

func (s *Service) ValidateRequest(req *Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	if req.StartTimeFrame.After(req.EndTimeFrame) ||
		req.StartTimeFrame.Before(time.Now().Add(4*24*time.Hour)) &&
			req.StartTimeFrame.Hour() >= 9 && req.StartTimeFrame.Hour() <= 23 &&
			req.StartTimeFrame.Before(time.Now()) {
		return errors.New("invalid time frame")
	}

	if _, err := s.repo.GetByOrderNumber(req.OrderNumber); err == nil {
		return errors.New("order already exists")
	}
	return nil
}

func (s *Service) CreateOrder(ctx context.Context, req *Request) error {
	if err := s.ValidateRequest(req); err != nil {
		return err
	}

	o := req.toOrder()

	if req.StartTimeFrame.Before(time.Now().Add(1 * time.Hour)) {
		s.SendToKafka(ctx, o)
	}

	return s.repo.Create(o)
}

func (s *Service) ProcessOrder(ctx context.Context) error {
	orders, err := s.repo.GetByTimeToDeliver(ctx)
	if err != nil {
		return err
	}
	if orders == nil {
		return errors.New("order not found")
	}

	for _, order := range orders {
		// Produce message to Kafka
		s.SendToKafka(ctx, &order)
	}

	return nil
}

func (s *Service) SendToKafka(ctx context.Context, order *entity.Order) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
		Value: func() []byte {
			data, err := json.Marshal(order)
			if err != nil {
				return nil
			}
			return data
		}(),
	}
	return s.kafkaProducer.Produce(message, nil)
}

// func (s *Service) ProcessDelivery(ctx context.Context, reqID string) error {
// 	req, err := s.repo.GetRequest(ctx, reqID)
// 	if err != nil {
// 		return err
// 	}

// 	if req.State != string(StateInit) {
// 		return errors.New("invalid state for processing")
// 	}

// 	req.State = string(StateIsFinding)
// 	if err := s.repo.UpdateRequest(ctx, req); err != nil {
// 		return err
// 	}

// 	// Simulate 3PL integration
// 	shipmentFound := s.fake3PLCall(req)
// 	if shipmentFound {
// 		req.State = string(StateFound)
// 	} else {
// 		req.State = string(StateNotFound)
// 	}
// 	return s.repo.UpdateRequest(ctx, req)
// }

// func (s *Service) fake3PLCall(req *Request) bool {
// 	// Simulate shipment search with success rate
// 	return time.Now().Unix()%2 == 0 // 50% success
// }
