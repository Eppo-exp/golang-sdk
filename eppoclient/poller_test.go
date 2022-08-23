package eppoclient

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_PollerPoll_InvokesCallbackUntilStoped(t *testing.T) {
	callCount := 0
	expected := 5

	var poller = NewPoller(1, func() {
		callCount++
	})
	poller.Start()
	time.Sleep(5 * time.Second)
	poller.Stop()

	assert.Equal(t, expected, callCount)
}

func Test_PollerPoll_StopsOnError(t *testing.T) {
	callCount := 0
	expected := 3

	var poller = NewPoller(1, func() {
		callCount++
		if callCount == 3 {
			panic("some_error")
		}
	})
	poller.Start()

	time.Sleep(5 * time.Second)
	assert.Equal(t, expected, callCount)
}
