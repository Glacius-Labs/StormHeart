package event

type EventSink interface {
	Handle(event DispatcherEvent)
}
