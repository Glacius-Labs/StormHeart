package event

import "sync"

type Dispatcher struct {
	mu    sync.Mutex
	sinks []EventSink
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		sinks: make([]EventSink, 0),
	}
}

func (d *Dispatcher) Register(sink EventSink) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.sinks = append(d.sinks, sink)
}

func (d *Dispatcher) Dispatch(event DispatchableEvent) {
	dispatcherEvent := event.ToDispatcherEvent()

	d.mu.Lock()
	defer d.mu.Unlock()

	for _, sink := range d.sinks {
		sink.Handle(dispatcherEvent)
	}
}
