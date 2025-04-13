package event

import (
	"fmt"
	"time"
)

const EventTypeHandler Type = "handler"

type HandlerEvent struct {
	handler   string
	event     Event
	err       error
	timestamp time.Time
}

func NewHandlerEvent(handler string, event Event, err error) *HandlerEvent {
	return &HandlerEvent{
		handler:   handler,
		event:     event,
		err:       err,
		timestamp: time.Now(),
	}
}

func (e *HandlerEvent) Type() Type {
	return EventTypeHandler
}

func (e *HandlerEvent) Message() string {
	if e.err != nil {
		return fmt.Sprintf("Handler '%s' encountered an error with event %v: %v", e.handler, e.event, e.err)
	}
	return fmt.Sprintf("Handler '%s' executed successfully with event: %v", e.handler, e.event)
}

func (e *HandlerEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *HandlerEvent) Error() error {
	return e.err
}
