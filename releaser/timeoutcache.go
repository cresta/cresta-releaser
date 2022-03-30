package releaser

import (
	"sync"
	"time"
)

type expireValues[V any] struct {
	value    V
	expireAt time.Time
}

type ExpireCache[K comparable, V any] struct {
	cache         map[K]expireValues[V]
	defaultExpiry time.Duration
	mu            sync.Mutex
}

func (e *ExpireCache[K, V]) Get(key K) (V, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if v, ok := e.cache[key]; ok {
		if v.expireAt.After(time.Now()) {
			return v.value, true
		}
		delete(e.cache, key)
	}
	var ret V
	return ret, false
}

func (e *ExpireCache[K, V]) Set(key K, value V) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cache == nil {
		e.cache = make(map[K]expireValues[V])
	}
	e.cache[key] = expireValues[V]{value, time.Now().Add(e.defaultExpiry)}
}
