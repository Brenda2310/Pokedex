package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]CacheEntry
	mu    sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	result := &Cache{
		cache: make(map[string]CacheEntry),
	}
	go result.reapLoop(interval)
	return result
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	result := CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	c.cache[key] = result
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cacheEn, ok := c.cache[key]
	if !ok {
		return []byte{}, false
	}
	res := cacheEn.val
	return res, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.cache {
			since := time.Since(entry.createdAt)
			if since > interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
