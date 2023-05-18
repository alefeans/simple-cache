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

// Expired returns true if Expiration is lower than current time
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

// New returns a new Cache using defaultExpiration as expiration
// time when adding an entry using SetDefault. The cleanupInterval
// is used to remove entries from the cache that are already expired.
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		entries:           make(map[string]Entry),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		stopCleanup:       make(chan bool),
	}
	go c.cleanupExpired()
	return c
}

// Get an entry from the cache. Returns the entry value or nil, and a bool
// indicating if it was found.
func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, found := c.entries[key]
	c.mu.RUnlock()
	if !found || e.Expired() {
		return nil, false
	}
	return e.Value, true
}

// Set adds an entry to the cache, replacing existing entry if has the same key.
// If expiration is lower than 1, the entry never expires.
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

// SetDefault adds an entry to the cache, using Cache.defaultExpiration as
// the expiration time.
func (c *Cache) SetDefault(key string, value any) {
	c.Set(key, value, c.defaultExpiration)
}

// SetNoExpire adds an entry to the cache that never expires.
func (c *Cache) SetNoExpire(key string, value any) {
	c.Set(key, value, time.Duration(NoExpiration))
}

// Delete entry from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.entries, key)
	c.mu.Unlock()
}

// Clear all entries from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	c.entries = make(map[string]Entry)
	c.mu.Unlock()
}

func (c *Cache) StopCleanup() {
	c.stopCleanup <- true
	close(c.stopCleanup)
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

// Close clear all entries from the Cache and stops the cleanup goroutine,
// gracefully freeing all used resources.
func (c *Cache) Close() {
	c.Clear()
	c.StopCleanup()
}

// Length return the number of entries in the cache, possibly including
// expired entries that weren't cleaned yet.
func (c *Cache) Length() int {
	c.mu.RLock()
	length := len(c.entries)
	c.mu.RUnlock()
	return length
}
