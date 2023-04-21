package cache

import (
	"time"
)

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

const (
	NoExpiration int64 = -1
)

func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]Entry),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	v, found := c.entries[key]
	if !found {
		return nil, false
	}
	if v.Expired() {
		return nil, false
	}
	return v.Value, true
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
