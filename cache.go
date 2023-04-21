package cache

import (
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
	entries map[string]Entry
}

func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]Entry),
	}
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
		c.SetNoExpire(key, value)
		return
	}
	c.entries[key] = Entry{Value: value, Expiration: time.Now().Add(expiration).UnixNano()}
}

func (c *Cache) SetNoExpire(key string, value any) {
	c.entries[key] = Entry{Value: value, Expiration: NoExpiration}
}

func (c *Cache) Delete(key string) {
	delete(c.entries, key)
}

