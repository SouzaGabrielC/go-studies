package events

import (
	"errors"
	"sync"
)

var EventHandlerAlreadyRegistered = errors.New("handler already registered")
var EventHandlerNotFoundForEventName = errors.New("handler not found for event name")
var EventNotFound = errors.New("event not found")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (e *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := e.handlers[eventName]; ok {
		for _, handlerRegistered := range e.handlers[eventName] {
			if handlerRegistered == handler {
				return EventHandlerAlreadyRegistered
			}
		}
	}

	e.handlers[eventName] = append(e.handlers[eventName], handler)
	return nil
}

func (e *EventDispatcher) Unregister(eventName string, handler EventHandlerInterface) error {
	if handlers, ok := e.handlers[eventName]; ok {
		for index, handlerRegistered := range handlers {
			if handlerRegistered == handler {
				e.handlers[eventName] = append(handlers[:index], handlers[index+1:]...)
				return nil
			}
		}

		return EventHandlerNotFoundForEventName
	}

	return EventNotFound
}

func (e *EventDispatcher) Dispatch(event EventInterface) {
	if handlers, ok := e.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)

			go func(h EventHandlerInterface) {
				h.Handle(event)
				wg.Done()
			}(handler)
		}
		wg.Wait()
	}
}

func (e *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	handlers, ok := e.handlers[eventName]
	if !ok {
		return false
	}

	for _, handlerRegistered := range handlers {
		if handlerRegistered == handler {
			return true
		}
	}

	return false
}

func (e *EventDispatcher) Clear() {
	clear(e.handlers)
}
