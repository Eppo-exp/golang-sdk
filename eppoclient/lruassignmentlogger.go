package eppoclient

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
)

type LruAssignmentLogger struct {
	cache *lru.TwoQueueCache[cacheKey, cacheValue]
	inner IAssignmentLogger
}

// We are only interested in whether a subject was ever a part of an
// assignment. We are not interested in the order of assignments or
// knowing the latest assignment. Therefore, both allocation and
// variation are part of the cacheKey.
type cacheKey struct {
	flag    string
	subject string
}
type cacheValue struct {
	allocation string
	variation  string
}

func NewLruAssignmentLogger(logger IAssignmentLogger, cacheSize int) (IAssignmentLogger, error) {
	cache, err := lru.New2Q[cacheKey, cacheValue](cacheSize)
	if err != nil {
		// err is only returned if `cacheSize` is invalid
		// (e.g., <0) which should normally never happen.
		return nil, fmt.Errorf("failed to create LRU cache: %w", err)
	}
	return &LruAssignmentLogger{
		cache: cache,
		inner: logger,
	}, nil
}

func (lal *LruAssignmentLogger) LogAssignment(event AssignmentEvent) {
	key := cacheKey{
		flag:    event.FeatureFlag,
		subject: event.Subject,
	}
	value := cacheValue{
		allocation: event.Allocation,
		variation:  event.Variation,
	}
	previousValue, recentlyLogged := lal.cache.Get(key)
	if !recentlyLogged || previousValue != value {
		lal.inner.LogAssignment(event)
		// Adding to cache after `LogAssignment` returned in
		// case it panics.
		lal.cache.Add(key, value)
	}
}
