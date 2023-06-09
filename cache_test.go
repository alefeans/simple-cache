package cache

import (
	"testing"
	"time"
)

func TestGetAndSetNoExpire(t *testing.T) {
	c := New(10*time.Millisecond, 20*time.Millisecond)

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

	c.SetNoExpire("replace", 2)
	c.SetNoExpire("replace", 3)
	v, found = c.Get("replace")
	if !found || v == nil || v != 3 {
		t.Errorf("Got value %v, found %t, but should be 3", v, found)
	}
}

func TestSet(t *testing.T) {
	c := New(10*time.Millisecond, 20*time.Millisecond)

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

	c.Set("replace set", 2, 5*time.Second)
	c.Set("replace set", 3, 5*time.Second)
	v, found = c.Get("replace set")
	if !found || v == nil || v != 3 {
		t.Errorf("Got value %v, found %t, but should be 3", v, found)
	}

	c.Set("expired", 1, 1*time.Nanosecond)
	<-time.After(2 * time.Nanosecond)
	v, found = c.Get("expired")
	if found || v == 1 {
		t.Errorf("Got value %v and found %t, but should be expired", v, found)
	}

	c.Set("replace timeout", 2, 10*time.Millisecond)
	c.Set("replace timeout", 3, 20*time.Millisecond)
	<-time.After(10 * time.Millisecond)
	v, found = c.Get("replace timeout")
	if !found || v == nil || v != 3 {
		t.Errorf("Got value %v, found %t, but should be 3 and not expired", v, found)
	}

	c.Set("replace expired", 3, 20*time.Millisecond)
	c.Set("replace expired", 2, 10*time.Millisecond)
	<-time.After(10 * time.Millisecond)
	v, found = c.Get("replace expired")
	if found || v == 3 {
		t.Errorf("Got value %v, found %t, but should be expired", v, found)
	}
}

func TestSetDefault(t *testing.T) {
	c := New(10*time.Millisecond, 20*time.Millisecond)

	c.SetDefault("k", "v")
	v, found := c.Get("k")
	if !found || v == nil || v != "v" {
		t.Errorf("Got value %v, found %t", v, found)
	}

	<-time.After(10 * time.Millisecond)
	v, found = c.Get("k")
	if found || v == "v" {
		t.Errorf("Got value %v, found %t, but should be expired", v, found)
	}
}

func TestDelete(t *testing.T) {
	c := New(10*time.Millisecond, 20*time.Millisecond)

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

	c.SetNoExpire("k1", "v1")
	c.SetNoExpire("k2", "v2")
	c.Delete("k1")
	v, found = c.Get("k2")
	if !found || v == nil || v != "v2" {
		t.Errorf("Got value %v, found %t", v, found)
	}

	if len(c.entries) != 1 {
		t.Error("Empty cache entries length is different from 1")
	}
}

func TestClear(t *testing.T) {
	c := New(10*time.Millisecond, 20*time.Millisecond)

	c.SetNoExpire("k1", "v1")
	c.SetNoExpire("k2", "v2")
	c.SetNoExpire("k3", "v3")
	c.Clear()
	v, found := c.Get("k1")
	if found || v != nil {
		t.Errorf("Got value %v, found %t, but should be nil", v, found)
	}

	if len(c.entries) != 0 {
		t.Error("Empty cache entries length is different from 0")
	}
}

func TestCleanupExpired(t *testing.T) {
	c := New(2*time.Millisecond, 5*time.Millisecond)

	c.SetNoExpire("k1", "v1")
	c.SetDefault("k2", "v2")
	c.Set("k3", "v3", 1*time.Second)
	c.Set("k4", "v4", 1*time.Millisecond)
	c.Set("k5", "v5", 5*time.Millisecond)
	<-time.After(6 * time.Millisecond)
	if len(c.entries) != 2 {
		t.Error("Cache entries length is different from 2")
	}
}

func TestStopCleanup(t *testing.T) {
	c := New(2*time.Millisecond, 5*time.Millisecond)

	c.StopCleanup()
	c.SetNoExpire("k1", "v1")
	c.Set("k2", "v2", 1*time.Second)
	c.Set("k3", "v3", 1*time.Millisecond)
	c.Set("k4", "v4", 5*time.Millisecond)
	<-time.After(6 * time.Millisecond)
	if len(c.entries) != 4 {
		t.Errorf("Cache entries length is different from 4")
	}
}

func TestCloseCache(t *testing.T) {
	c := New(2*time.Millisecond, 5*time.Millisecond)
	c.SetNoExpire("k1", "v1")
	c.Set("k2", "v2", 1*time.Second)
	c.Set("k3", "v3", 1*time.Millisecond)
	c.Set("k4", "v4", 5*time.Millisecond)
	c.Close()
	if len(c.entries) > 0 {
		t.Errorf("Closed cache entries length is different from 0")
	}
	select {
	case <-c.stopCleanup:
	default:
		t.Error("Channel is not closed")
	}
}

func TestLength(t *testing.T) {
	c := New(2*time.Millisecond, 5*time.Millisecond)

	c.SetNoExpire("k1", "v1")
	c.SetNoExpire("k2", "v2")
	c.SetNoExpire("k3", "v3")
	c.SetNoExpire("k4", "v4")
	if c.Length() != 4 {
		t.Errorf("Cache length is different from 4")
	}
}

func BenchmarkCacheGet(b *testing.B) {
	b.StopTimer()
	c := New(2*time.Second, 5*time.Second)
	c.SetNoExpire("k", "v")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Get("k")
	}
}

func BenchmarkCacheSet(b *testing.B) {
	b.StopTimer()
	c := New(2*time.Second, 5*time.Second)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Set("k", "v", 3*time.Second)
	}
}

func BenchmarkCacheSetDefault(b *testing.B) {
	b.StopTimer()
	c := New(2*time.Second, 5*time.Second)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.SetDefault("k", "v")
	}
}

func BenchmarkCacheSetNoExpire(b *testing.B) {
	b.StopTimer()
	c := New(2*time.Second, 5*time.Second)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.SetNoExpire("k", "v")
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	b.StopTimer()
	c := New(2*time.Second, 5*time.Second)
	c.SetNoExpire("k", "v")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Delete("k")
	}
}
