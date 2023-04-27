package cache

import "time"

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
	e, found := c.entries[key]
	if !found || e.Expired() {
		return nil, false
	}
	return e.Value, true
}

func (c *Cache) Set(key string, value any, expiration time.Duration) {
	if expiration < 1 {
		c.entries[key] = Entry{Value: value, Expiration: NoExpiration}
		return
	}
	c.entries[key] = Entry{Value: value, Expiration: time.Now().Add(expiration).UnixNano()}
}

func (c *Cache) SetDefault(key string, value any) {
	c.Set(key, value, c.defaultExpiration)
}

func (c *Cache) SetNoExpire(key string, value any) {
	c.Set(key, value, time.Duration(NoExpiration))
}

func (c *Cache) Delete(key string) {
	delete(c.entries, key)
}

func (c *Cache) Clear() {
	c.entries = make(map[string]Entry)
}

func (c *Cache) StopCleanup() {
	c.stopCleanup <- true
}

func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for k, e := range c.entries {
				if e.Expired() {
					c.Delete(k)
				}
			}
		case <-c.stopCleanup:
			return
		}
	}
}
