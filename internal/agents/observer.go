package agents

import "github.com/sanjbh/mailforge/internal/events"

type Observable struct {
	Observers []events.Observer
}

func (obs *Observable) Register(observer events.Observer) {
	obs.Observers = append(obs.Observers, observer)
}

func (obs *Observable) NotifyAll(event events.AgentEvent) {
	for _, observer := range obs.Observers {
		observer.Notify(event)
	}
}
