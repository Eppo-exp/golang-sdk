package eppoclient

import (
	"testing"
	"time"
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

	if callCount != expected {
		t.Errorf("\"Poller\" FAILED, expected -> %v, got -> %v", expected, callCount)
	} else {
		t.Logf("\"Poller\" SUCCEDED, expected -> %v, got -> %v", expected, callCount)
	}
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

	if callCount != expected {
		t.Errorf("\"Poller\" FAILED, expected -> %v, got -> %v", expected, callCount)
	} else {
		t.Logf("\"Poller\" SUCCEDED, expected -> %v, got -> %v", expected, callCount)
	}
}
