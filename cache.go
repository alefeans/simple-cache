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
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		entries:           make(map[string]Entry),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
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
