package mock

import (
	"context"
	"sync"

	"slices"

	"github.com/glacius-labs/StormHeart/internal/core/event"
)

type MockHandler struct {
	mu     sync.Mutex
	events []event.EventType
}

func NewMockHandler() *MockHandler {
	return &MockHandler{
		events: make([]event.EventType, 0),
	}
}

func (h *MockHandler) Handle(ctx context.Context, e event.EventType) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = append(h.events, e)
}

func (h *MockHandler) Events() []event.EventType {
	h.mu.Lock()
	defer h.mu.Unlock()

	return slices.Clone(h.events)
}

func (h *MockHandler) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = nil
}
