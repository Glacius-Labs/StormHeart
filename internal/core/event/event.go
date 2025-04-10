package event

import "time"

type Event interface {
	Message() string
	Type() Type
	Error() error
	Timestamp() time.Time
}
