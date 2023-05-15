package cache

import (
	"sync"
	"time"
)

const NoExpiration int64 = -1

type Entry struct {
	Expiration int64
	Value      any
}

func (e *Entry) Expired() bool {
	if e.Expiration == NoExpiration {
		return false
	}
	return time.Now().UnixNano() > e.Expiration
}

type Cache struct {
	mu                sync.RWMutex
	entries           map[string]Entry
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	stopCleanup       chan bool
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		entries:           make(map[string]Entry),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		stopCleanup:       make(chan bool),
	}
	go c.cleanupExpired()
	return c
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, found := c.entries[key]
	c.mu.RUnlock()
	if !found || e.Expired() {
		return nil, false
	}
	return e.Value, true
}

func (c *Cache) Set(key string, value any, expiration time.Duration) {
	if expiration < 1 {
		c.set(key, value, NoExpiration)
		return
	}
	c.set(key, value, time.Now().Add(expiration).UnixNano())
}

func (c *Cache) set(key string, value any, expiration int64) {
	c.mu.Lock()
	c.entries[key] = Entry{Value: value, Expiration: expiration}
	c.mu.Unlock()
}

func (c *Cache) SetDefault(key string, value any) {
	c.Set(key, value, c.defaultExpiration)
}

func (c *Cache) SetNoExpire(key string, value any) {
	c.Set(key, value, time.Duration(NoExpiration))
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.entries, key)
	c.mu.Unlock()
}

func (c *Cache) Clear() {
	c.mu.Lock()
	c.entries = make(map[string]Entry)
	c.mu.Unlock()
}

func (c *Cache) StopCleanup() {
	c.stopCleanup <- true
}

func (c *Cache) deleteExpired() {
	c.mu.Lock()
	for k, e := range c.entries {
		if e.Expired() {
			delete(c.entries, k)
		}
	}
	c.mu.Unlock()
}

func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCleanup:
			return
		}
	}
}
