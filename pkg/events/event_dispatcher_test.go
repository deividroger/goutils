package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (ed *TestEvent) GetName() string {
	return ed.Name
}

func (ed *TestEvent) GetPayload() interface{} {
	return ed.Payload
}

func (ed *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(e EventInterface, wg *sync.WaitGroup) {

}

type EventDispatcherTestSuite struct {
	suite.Suite

	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	EventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.EventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandler{
		ID: 1,
	}
	suite.handler2 = TestEventHandler{
		ID: 2,
	}
	suite.handler3 = TestEventHandler{
		ID: 3,
	}

	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}

}

func (suite *EventDispatcherTestSuite) TestEventDispacher_Register() {
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))

	err = suite.EventDispatcher.Register(suite.event.Name, &suite.handler2)
	suite.Nil(err)

	suite.Equal(2, len(suite.EventDispatcher.handlers[suite.event.Name]))

	assert.Equal(suite.T(), &suite.handler, suite.EventDispatcher.handlers[suite.event.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.EventDispatcher.handlers[suite.event.GetName()][1])
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(e EventInterface, wg *sync.WaitGroup) {
	m.Called(e)
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatch_Dispatch() {
	eh := &MockHandler{}
	eh.On("Handle", &suite.event)

	eh2 := &MockHandler{}
	eh2.On("Handle", &suite.event)

	suite.EventDispatcher.Register(suite.event.Name, eh)
	suite.EventDispatcher.Register(suite.event.Name, eh2)

	suite.EventDispatcher.Dispatch(&suite.event)

	eh.AssertExpectations(suite.T())
	eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
	eh2.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *EventDispatcherTestSuite) TestEvent_Has() {
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)

	suite.True(suite.EventDispatcher.Has(suite.event.Name, &suite.handler))
	suite.False(suite.EventDispatcher.Has(suite.event.Name, &suite.handler2))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))

	suite.EventDispatcher.Clear()

	suite.Equal(0, len(suite.EventDispatcher.handlers[suite.event.Name]))
}

func (suite *EventDispatcherTestSuite) TestEventDispacher_Remove() {
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))

	suite.EventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.Equal(0, len(suite.EventDispatcher.handlers[suite.event.Name]))
}

func (suite *EventDispatcherTestSuite) TestEventDispacher_Register_ErrHandlerAlreadyRegistered() {
	err := suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))

	err = suite.EventDispatcher.Register(suite.event.Name, &suite.handler)
	suite.Equal(ErrHandlerAlreadyRegistered, err)

	suite.Equal(1, len(suite.EventDispatcher.handlers[suite.event.Name]))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
