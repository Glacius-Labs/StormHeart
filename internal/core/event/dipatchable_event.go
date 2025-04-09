package event

type DispatchableEvent interface {
	ToDispatcherEvent() DispatcherEvent
}
