package cache

import (
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	c := NewCache()
	v, found := c.Get("new")
	if found || v != nil {
		t.Errorf("Found %v value with empty cache", v)
	}

	c.SetNoExpire("k", "v")
	v, found = c.Get("k")
	if !found || v == nil || v != "v" {
		t.Errorf("Got value %v, found %t", v, found)
	}

	c.SetNoExpire("x", 1)
	if len(c.entries) != 2 {
		t.Error("Cache entries length is different from 2")
	}

	x := 1
	c.SetNoExpire("pointer", &x)
	p, _ := c.Get("pointer")
	if pp, ok := p.(*int); ok {
		if x != *pp {
			t.Errorf("Got value %v different from %v", *pp, x)
		}
	} else {
		t.Errorf("Expected a pointer to int, but got %T", p)
	}
}

func TestExpiration(t *testing.T) {
	c := NewCache()

	c.Set("negative", "expiration", -1)
	v, found := c.Get("negative")
	if !found || v == nil || v != "expiration" {
		t.Errorf("Got value %v, found %t", v, found)
	}

	c.Set("zero", "expiration", 0)
	v, found = c.Get("zero")
	if !found || v == nil || v != "expiration" {
		t.Errorf("Got value %v, found %t", v, found)
	}

	c.Set("k", "v", 5*time.Second)
	v, found = c.Get("k")
	if !found || v == nil || v != "v" {
		t.Errorf("Got value %v, found %t", v, found)
	}

	c.Set("expired", 1, 1*time.Nanosecond)
	<-time.After(2 * time.Nanosecond)
	v, found = c.Get("expired")
	if found || v == 1 {
		t.Errorf("Got value %v and found %t, but should be expired", v, found)
	}
}

func TestDelete(t *testing.T) {
	c := NewCache()

	c.Delete("x") // delete from empty cache is ok

	c.SetNoExpire("k", "v")
	c.Delete("k")
	v, found := c.Get("k")
	if found || v != nil {
		t.Errorf("Got value %v, found %t, but should be nil", v, found)
	}

	if len(c.entries) != 0 {
		t.Error("Empty cache entries length is different from 0")
	}
}
