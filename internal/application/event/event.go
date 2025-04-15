package event

import "time"

type Event interface {
	Message() string
	Type() EventType
	Error() error
	Timestamp() time.Time
}
