package agents

type EventType int

const (
	EventProgress EventType = iota
	EventSuccess
	EventError
)

type Observer interface {
	Notify(event AgentEvent)
}

type AgentEvent struct {
	Type    EventType
	Payload int
}

type Observable struct {
	Observers []Observer
}

func (obs *Observable) Register(observer Observer) {
	obs.Observers = append(obs.Observers, observer)
}

func (obs *Observable) NotifyAll(event AgentEvent) {
	for _, observer := range obs.Observers {
		observer.Notify(event)
	}
}
