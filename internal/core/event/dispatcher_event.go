package event

import "time"

type DispatcherEvent struct {
	Message   string
	Type      DispatcherEventType
	Error     error
	Timestamp time.Time
}
