package eppoclient

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
)

type LruBanditLogger struct {
	cache *lru.TwoQueueCache[lruBanditKey, lruBanditValue]
	inner IAssignmentLogger
}

type lruBanditKey struct {
	flagKey    string
	subjectKey string
}
type lruBanditValue struct {
	banditKey string
	actionKey string
}

func NewLruBanditLogger(logger IAssignmentLogger, cacheSize int) (IAssignmentLogger, error) {
	cache, err := lru.New2Q[lruBanditKey, lruBanditValue](cacheSize)
	if err != nil {
		// err is only returned if `cacheSize` is invalid
		// (e.g., <0) which should normally never happen.
		return nil, fmt.Errorf("failed to create LRU cache: %w", err)
	}
	return &LruBanditLogger{
		cache: cache,
		inner: logger,
	}, nil
}

func (logger *LruBanditLogger) LogAssignment(event AssignmentEvent) {
	logger.inner.LogAssignment(event)
}

func (logger *LruBanditLogger) LogBanditAction(event BanditEvent) {
	inner, ok := logger.inner.(BanditActionLogger)
	if !ok {
		return
	}

	key := lruBanditKey{
		flagKey:    event.FlagKey,
		subjectKey: event.Subject,
	}
	value := lruBanditValue{
		banditKey: event.BanditKey,
		actionKey: event.Action,
	}
	previousValue, recentlyLogged := logger.cache.Get(key)
	if !recentlyLogged || previousValue != value {
		inner.LogBanditAction(event)
		// Adding to cache after `LogBanditAction` returned in
		// case it panics.
		logger.cache.Add(key, value)
	}
}
