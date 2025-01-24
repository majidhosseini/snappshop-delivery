package delivery

import (
	"context"
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

func NewService(repo repository.OrderRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ValidateRequest(req *Request) error {
	if req.OrderNumber == "" || req.FromLoc == [2]float64{} || req.ToLoc == [2]float64{} {
		return errors.New("invalid request")
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

	return s.repo.Create(req.toOrder())
}

func (s *Service) ProcessOrder(ctx context.Context, orderID uint64) error {
	order, err := s.orderRepository.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	if order.Status != entity.StatusCreated {
		return errors.New("invalid state for processing")
	}

	order.Status = entity.StatusInProgress
	order.UpdatedAt = time.Now()
	if err := s.orderRepository.Update(order); err != nil {
		return err
	}

	// Simulate 3PL integration
	shipmentFound := s.fake3PLCall(order)
	if shipmentFound {
		order.Status = entity.StatusCompleted
	} else {
		order.Status = entity.StatusCanceled
	}
	order.UpdatedAt = time.Now()
	if err := s.orderRepository.Update(order); err != nil {
		return err
	}

	// Produce message to Kafka
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
		Value:          []byte("Order processed: " + order.OrderNumber),
	}
	if err := s.kafkaProducer.Produce(message, nil); err != nil {
		return err
	}

	return nil
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
