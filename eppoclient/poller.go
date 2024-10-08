package eppoclient

import (
	"time"
)

type poller struct {
	interval          time.Duration
	callback          func()
	isStopped         bool `default:"false"`
	applicationLogger ApplicationLogger
}

func newPoller(interval time.Duration, callback func(), applicationLogger ApplicationLogger) *poller {
	var pl = &poller{}

	pl.interval = interval
	pl.callback = callback
	pl.applicationLogger = applicationLogger

	return pl
}

func (p *poller) Start() {
	if p.applicationLogger != nil {
		p.applicationLogger.Info("Poller start")
	}
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
	if p.applicationLogger != nil {
		p.applicationLogger.Info("Poller stopped")
	}
	p.isStopped = true
}
