package event

import (
	"context"
	"sync"
)

type Dispatcher struct {
	mu       sync.Mutex
	handlers []Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make([]Handler, 0),
	}
}

func (d *Dispatcher) Register(handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.handlers = append(d.handlers, handler)
}

func (d *Dispatcher) Dispatch(ctx context.Context, event Event) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, handler := range d.handlers {
		handler.Handle(ctx, event)
	}
}
