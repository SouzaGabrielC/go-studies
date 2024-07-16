package events

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type testEvent struct {
	name      string
	payload   interface{}
	timestamp time.Time
}

func (e *testEvent) GetName() string {
	return e.name
}

func (e *testEvent) GetPayload() interface{} {
	return e.payload
}

func (e *testEvent) GetDateTime() time.Time {
	return e.timestamp
}

type testHandler struct {
	ID int
}

func (h *testHandler) Handle(event EventInterface) {}

type EventDispatcherTestSuite struct {
	suite.Suite
	event      EventInterface
	event2     EventInterface
	handler    EventHandlerInterface
	handler2   EventHandlerInterface
	handler3   EventHandlerInterface
	dispatcher *EventDispatcher
}

func (s *EventDispatcherTestSuite) SetupTest() {
	s.dispatcher = NewEventDispatcher()
	s.handler = &testHandler{
		ID: 1,
	}
	s.handler2 = &testHandler{
		ID: 2,
	}
	s.handler3 = &testHandler{
		ID: 3,
	}
	s.event = &testEvent{
		name:      "testEvent",
		timestamp: time.Now(),
		payload: struct {
			message string
		}{
			message: "Test message",
		},
	}
	s.event2 = &testEvent{
		name:      "testEvent 2",
		timestamp: time.Now(),
		payload: struct {
			message string
		}{
			message: "Test message 2",
		},
	}
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := s.dispatcher.Register(s.event.GetName(), s.handler)
	s.Nil(err)
	s.Equal(1, len(s.dispatcher.handlers[s.event.GetName()]))

	err = s.dispatcher.Register(s.event.GetName(), s.handler2)
	s.Nil(err)
	s.Equal(2, len(s.dispatcher.handlers[s.event.GetName()]))
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := s.dispatcher.Register(s.event.GetName(), s.handler)
	s.Nil(err)

	err = s.dispatcher.Register(s.event.GetName(), s.handler)
	s.ErrorIs(err, EventHandlerAlreadyRegistered)
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	// Register event 1 handler
	err := s.dispatcher.Register(s.event.GetName(), s.handler)
	s.Nil(err)

	// Register event 2 handler
	err = s.dispatcher.Register(s.event2.GetName(), s.handler2)
	s.Nil(err)

	// Check if events are registered
	_, ok := s.dispatcher.handlers[s.event.GetName()]
	s.True(ok)

	_, ok = s.dispatcher.handlers[s.event2.GetName()]
	s.True(ok)

	s.Equal(2, len(s.dispatcher.handlers))

	// Run clear
	s.dispatcher.Clear()

	// Make sure we don't have any event registered
	_, ok = s.dispatcher.handlers[s.event.GetName()]
	s.False(ok)

	_, ok = s.dispatcher.handlers[s.event2.GetName()]
	s.False(ok)

	s.Equal(len(s.dispatcher.handlers), 0)
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	// Register event 1 handler
	err := s.dispatcher.Register(s.event.GetName(), s.handler)
	s.Nil(err)

	err = s.dispatcher.Register(s.event.GetName(), s.handler2)
	s.Nil(err)

	hasEventHandler := s.dispatcher.Has(s.event.GetName(), s.handler)
	s.True(hasEventHandler)

	hasEventHandler = s.dispatcher.Has(s.event.GetName(), s.handler2)
	s.True(hasEventHandler)

	hasEventHandler = s.dispatcher.Has(s.event2.GetName(), s.handler2)
	s.False(hasEventHandler)
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface) {
	m.Called(event)
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	handler := &MockHandler{}
	handler.On("Handle", s.event)

	err := s.dispatcher.Register(s.event.GetName(), handler)
	s.Nil(err)

	s.dispatcher.Dispatch(s.event)

	handler.AssertExpectations(s.T())
	handler.AssertNumberOfCalls(s.T(), "Handle", 1)
}

func (s *EventDispatcherTestSuite) TestEventDispatcher_Unregister() {
	// Register event 1 handler
	err := s.dispatcher.Register(s.event.GetName(), s.handler)
	s.Nil(err)

	err = s.dispatcher.Register(s.event.GetName(), s.handler2)
	s.Nil(err)

	err = s.dispatcher.Register(s.event.GetName(), s.handler3)
	s.Nil(err)

	s.Equal(3, len(s.dispatcher.handlers[s.event.GetName()]))

	err = s.dispatcher.Unregister(s.event.GetName(), s.handler2)
	s.Nil(err)

	s.Equal(2, len(s.dispatcher.handlers[s.event.GetName()]))

	hasHandler := s.dispatcher.Has(s.event.GetName(), s.handler2)
	s.False(hasHandler)
}

func TestSuiteEventDispatcher(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
