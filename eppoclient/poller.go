package eppoclient

import (
	"fmt"
	"time"
)

type Poller struct {
	interval  int `default:"10"`
	callback  func()
	isStopped bool `default:"false"`
}

func NewPoller(interval int, callback func()) *Poller {
	var poller = &Poller{}

	poller.interval = interval
	poller.callback = callback

	return poller
}

func (p *Poller) Start() {
	fmt.Println("Poller start")

	go p.poll()
}

func (p *Poller) poll() {
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

func (p *Poller) Stop() {
	fmt.Println("Poller stopped")
	p.isStopped = true
}
