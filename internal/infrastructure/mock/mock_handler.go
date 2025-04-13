package mock

import (
	"context"
	"sync"

	"slices"

	"github.com/glacius-labs/StormHeart/internal/core/event"
)

type MockHandler struct {
	mu     sync.Mutex
	events []event.Event
}

func NewMockHandler() *MockHandler {
	return &MockHandler{
		events: make([]event.Event, 0),
	}
}

func (h *MockHandler) Name() string {
	return "MockHandler"
}

func (h *MockHandler) Handle(ctx context.Context, e event.Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = append(h.events, e)
	return nil
}

func (h *MockHandler) Events() []event.Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	return slices.Clone(h.events)
}

func (h *MockHandler) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = nil
}
