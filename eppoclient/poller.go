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

func (p *Poller) New(interval int, callback func()) {
	p.interval = interval
	p.callback = callback
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
