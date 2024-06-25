package eppoclient

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CallbackMock struct {
	mock.Mock
}

func (m *CallbackMock) CallbackFn() {
	m.Called()
}

func Test_PollerPoll_InvokesCallbackUntilStoped(t *testing.T) {
	expected := 5

	callbackMock := CallbackMock{}
	callbackMock.On("CallbackFn").Return()

	var poller = newPoller(1*time.Second, callbackMock.CallbackFn)
	poller.Start()
	time.Sleep(5 * time.Second + 500 * time.Millisecond) // half second buffer to allow polling thread to execute
	poller.Stop()

	callbackMock.AssertNumberOfCalls(t, "CallbackFn", expected)
}

func Test_PollerPoll_StopsOnError(t *testing.T) {
	callCount := 0
	expected := 3

	var poller = newPoller(1*time.Second, func() {
		callCount++
		if callCount == 3 {
			panic("some_error")
		}
	})
	poller.Start()

	time.Sleep(5 * time.Second)
	assert.Equal(t, expected, callCount)
}

func Test_PollerPoll_ManualStop(t *testing.T) {
	expected := 3

	callbackMock := CallbackMock{}
	callbackMock.On("CallbackFn").Return()

	var poller = newPoller(1*time.Second, callbackMock.CallbackFn)
	poller.Start()

	time.Sleep(2500 * time.Millisecond)

	poller.Stop()

	time.Sleep(2 * time.Second)
	callbackMock.AssertNumberOfCalls(t, "CallbackFn", expected)
}
