package main

import (
	"sync"
	"time"
)

type cacheItem struct {
	value   string
	expires int64
}

func (i *cacheItem) isExpired() bool {
	return i.expires < time.Now().Unix()
}

type cache struct {
	mu    sync.RWMutex
	ttl   int64
	items map[string]cacheItem
}

func NewCache(ttl int64) cache {
	return cache{
		ttl:   ttl,
		items: make(map[string]cacheItem),
	}
}

func (c *cache) get(k string) (string, bool) {
	c.mu.RLock()
	item, found := c.items[k]
	c.mu.RUnlock()

	if !found || item.isExpired() {
		return "", false
	}

	return item.value, true
}

func (c *cache) set(k string, v string) {
	item := cacheItem{
		value:   v,
		expires: time.Now().Unix() + c.ttl,
	}

	c.mu.Lock()
	c.items[k] = item
	c.mu.Unlock()
}
