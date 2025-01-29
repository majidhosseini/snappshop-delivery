package delivery

import (
	"context"
	"fmt"
)

type MockRepository struct {
	requests map[string]*Request
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		requests: make(map[string]*Request),
	}
}

func (m *MockRepository) SaveRequest(ctx context.Context, req *Request) error {
	m.requests[req.OrderNumber] = req
	return nil
}

func (m *MockRepository) GetRequest(ctx context.Context, reqID string) (*Request, error) {
	req, exists := m.requests[reqID]
	if !exists {
		return nil, fmt.Errorf("request %s not found", reqID)
	}
	return req, nil
}

func (m *MockRepository) UpdateRequest(ctx context.Context, req *Request) error {
	m.requests[req.OrderNumber] = req
	return nil
}
