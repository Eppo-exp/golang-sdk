package eppoclient

import (
	"fmt"
	"time"
)

type poller struct {
	interval  time.Duration
	callback  func()
	isStopped bool `default:"false"`
}

func newPoller(interval time.Duration, callback func()) *poller {
	var pl = &poller{}

	pl.interval = interval
	pl.callback = callback

	return pl
}

func (p *poller) Start() {
	fmt.Println("Poller start")

	go p.poll()
}

func (p *poller) poll() {
	defer func() {
		if err := recover(); err != nil {
			p.Stop()
		}
	}()

	for {
		if p.isStopped {
			break
		}
		p.callback()
		time.Sleep(p.interval)
	}
}

func (p *poller) Stop() {
	fmt.Println("Poller stopped")
	p.isStopped = true
}
