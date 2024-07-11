package eppoclient

import (
	"time"

	"github.com/Eppo-exp/golang-sdk/v5/eppoclient/applicationlogger"
)

type poller struct {
	interval          time.Duration
	callback          func()
	isStopped         bool `default:"false"`
	applicationLogger applicationlogger.Logger
}

func newPoller(interval time.Duration, callback func(), applicationLogger ...applicationlogger.Logger) *poller {
	var pl = &poller{}

	pl.interval = interval
	pl.callback = callback
	if len(applicationLogger) > 0 {
		pl.applicationLogger = applicationLogger[0]
	}

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
