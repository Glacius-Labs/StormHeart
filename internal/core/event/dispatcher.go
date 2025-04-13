package event

import (
	"context"
	"slices"
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
	handlers := slices.Clone(d.handlers)
	d.mu.Unlock()

	for _, handler := range handlers {
		go func(h Handler) {
			err := h.Handle(ctx, event)

			if _, isHandlerEvent := event.(*HandlerEvent); isHandlerEvent {
				return
			}

			handlerEvent := NewHandlerEvent(h.Name(), event, err)
			d.Dispatch(ctx, handlerEvent)
		}(handler)
	}
}
