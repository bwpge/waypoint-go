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
	items map[string]cacheItem
}

func NewCache() cache {
	return cache{
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
		expires: time.Now().Unix() + 5,
	}

	c.mu.Lock()
	c.items[k] = item
	c.mu.Unlock()
}
