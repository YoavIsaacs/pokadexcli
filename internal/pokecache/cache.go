package pokecache

import (
	"errors"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	Val       []byte
}

type Cache struct {
	Cache  map[string]cacheEntry
	lock   sync.Mutex
	expire time.Duration
}

func NewCache(expirePeriod time.Duration) *Cache {
	c := &Cache{
		Cache:  make(map[string]cacheEntry),
		expire: expirePeriod,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	entry := cacheEntry{
		createdAt: time.Now(),
		Val:       val,
	}
	c.Cache[key] = entry
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	data, ok := c.Cache[key]

	if !ok {
		return []byte{}, errors.New("item not in cache...\n")
	}
	return data.Val, nil
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.expire)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.reap()
	}
}

func (c *Cache) reap() {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	for key, entry := range c.Cache {
		if now.Sub(entry.createdAt) > c.expire {
			delete(c.Cache, key)
		}
	}
}
