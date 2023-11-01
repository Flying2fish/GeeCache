package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu sync.Mutex
	ca *lru.Cache

	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ca == nil {
		c.ca = lru.New(c.cacheBytes, nil)
	}

	c.ca.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ca == nil {
		return
	}

	if v, ok := c.ca.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
