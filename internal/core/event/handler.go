package event

import "context"

type Handler interface {
	Name() string
	Handle(ctx context.Context, event Event) error
}
