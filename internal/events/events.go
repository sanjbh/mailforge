package events

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
