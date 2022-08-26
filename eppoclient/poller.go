package eppoclient

import (
	"fmt"
	"time"
)

type poller struct {
	interval  int `default:"10"`
	callback  func()
	isStopped bool `default:"false"`
}

func newPoller(interval int, callback func()) *poller {
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
		time.Sleep(time.Duration(p.interval) * time.Second)
	}
}

func (p *poller) Stop() {
	fmt.Println("Poller stopped")
	p.isStopped = true
}
