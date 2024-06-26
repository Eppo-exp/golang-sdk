package eppoclient

import (
	"fmt"
	"sync"
	"time"
)

type poller struct {
	interval time.Duration
	callback func()
	stopChan chan struct{}
	stopOnce sync.Once
}

const defaultPollInterval = 10 * time.Second

func newPoller(interval time.Duration, callback func()) *poller {
	if interval == 0 {
		interval = defaultPollInterval
	}
	return &poller{
		interval: interval,
		callback: callback,
		stopChan: make(chan struct{}),
	}
}

func (p *poller) Start() {
	fmt.Println("Poller start")

	go p.poll()
}

func (p *poller) poll() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recovered from panic: %v\n", err)
		}
		p.Stop()
	}()

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopChan:
			return
		case <-ticker.C:
			p.callback()
		}
	}
}

func (p *poller) Stop() {
	p.stopOnce.Do(func() {
		fmt.Println("Poller stopped")
		close(p.stopChan)
	})
}
