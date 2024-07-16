package events

import "time"

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
}

type EventHandlerInterface interface {
	Handle(event EventInterface)
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Unregister(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface)
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}
