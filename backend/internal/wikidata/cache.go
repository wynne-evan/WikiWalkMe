package wikidata

import (
	"fmt"
	"sync"
	"time"
)

type cacheItem struct {
	targets   []Target
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.Mutex
	items map[string]cacheItem
	ttl   time.Duration
}

func NewMemoryCache(ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		items: make(map[string]cacheItem),
		ttl:   ttl,
	}
}

// Use rounded grid key
func (c *MemoryCache) generateKey(lat, lon float64) string {
	return fmt.Sprintf("%.3f, %.3f", lat, lon)
}

// Check the cache for unexpired targets
func (c *MemoryCache) Get(lat, lon float64) ([]Target, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.generateKey(lat, lon)
	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		delete(c.items, key)
		return nil, false
	}

	return item.targets, true
}

// Save targets to cache
func (c *MemoryCache) Set(lat, lon float64, targets []Target) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.generateKey(lat, lon)
	c.items[key] = cacheItem{
		targets:   targets,
		expiresAt: time.Now().Add(c.ttl),
	}
}
