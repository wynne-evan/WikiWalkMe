package routing

import (
	"fmt"
	"sync"
	"time"
)

type routeCacheItem struct {
	response  RouteResponse
	expiresAt time.Time
}

type RouteCache struct {
	mu    sync.Mutex
	items map[string]routeCacheItem
	ttl   time.Duration
}

func NewRouteCache(ttl time.Duration) *RouteCache {
	return &RouteCache{
		items: make(map[string]routeCacheItem),
		ttl:   ttl,
	}
}

func (c *RouteCache) generateKey(startLat, startLon, endLat, endLon, maxMinutes float64) string {
	return fmt.Sprintf("%.3f,%.3f-%.3f.%.3f-%.0f", startLat, startLon, endLat, endLon, maxMinutes)
}

func (c *RouteCache) Get(startLat, startLon, endLat, endLon, maxMinutes float64) (RouteResponse, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.generateKey(startLat, startLon, endLat, endLon, maxMinutes)
	item, found := c.items[key]
	if !found {
		return RouteResponse{}, false
	}

	if time.Now().After(item.expiresAt) {
		delete(c.items, key)
		return RouteResponse{}, false
	}

	return item.response, true
}

func (c *RouteCache) Set(startLat, startLon, endLat, endLon, maxMinutes float64, response RouteResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.generateKey(startLat, startLon, endLat, endLon, maxMinutes)
	c.items[key] = routeCacheItem{
		response:  response,
		expiresAt: time.Now().Add(c.ttl),
	}
}
