package delivery

import (
    "context"
    "errors"
    "time"
)

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) ValidateRequest(req *Request) error {
    if req.OrderID == "" || req.FromLoc == [2]float64{} || req.ToLoc == [2]float64{} {
        return errors.New("invalid request")
    }
    return nil
}

func (s *Service) CreateRequest(ctx context.Context, req *Request) error {
    if err := s.ValidateRequest(req); err != nil {
        return err
    }
    req.State = string(StateInit)
    return s.repo.SaveRequest(ctx, req)
}

func (s *Service) ProcessShipment(ctx context.Context, reqID string) error {
    req, err := s.repo.GetRequest(ctx, reqID)
    if err != nil {
        return err
    }

    if req.State != string(StateInit) {
        return errors.New("invalid state for processing")
    }

    req.State = string(StateIsFinding)
    if err := s.repo.UpdateRequest(ctx, req); err != nil {
        return err
    }

    // Simulate 3PL integration
    shipmentFound := s.fake3PLCall(req)
    if shipmentFound {
        req.State = string(StateFound)
    } else {
        req.State = string(StateNotFound)
    }
    return s.repo.UpdateRequest(ctx, req)
}

func (s *Service) fake3PLCall(req *Request) bool {
    // Simulate shipment search with success rate
    return time.Now().Unix()%2 == 0 // 50% success
}
