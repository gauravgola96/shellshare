package storage

import (
	"errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var (
	Cache       *Cache_
	ErrNilCache = errors.New("nil cache")
)

type valueItem struct {
	Value      string
	startTime  time.Time
	lastAccess int64
	Expires    int64
}

// Cache_ is a struct for caching.
type Cache_ struct {
	items map[string]*valueItem
	mu    sync.Mutex
}

func InitializeCache() error {
	if Cache == nil {
		subLogger := log.With().Str("module", "storage.cache").Logger()
		Cache = &Cache_{items: make(map[string]*valueItem)}
		go func() {
			t := time.NewTicker(time.Second)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					Cache.mu.Lock()
					for k, v := range Cache.items {
						if time.Now().UnixNano()-v.lastAccess > v.Expires {
							subLogger.Info().Msgf("%v has expires at %s", Cache.items, time.Now().String())
							delete(Cache.items, k)
						}
					}
					Cache.mu.Unlock()
				}
			}
		}()
	}
	return nil
}

// Get gets a value from a cache.
func (c *Cache_) Get(key string) (string, time.Time, error) {
	c.mu.Lock()
	if v, ok := c.items[key]; ok {
		v.lastAccess = time.Now().UnixNano()
		return v.Value, v.startTime, nil
	}
	c.mu.Unlock()
	return "", time.Time{}, ErrNilCache
}

// Put puts a value to a cache. If a key and value exists, overwrite it.
func (c *Cache_) Put(key string, value string, expire time.Duration) {
	c.mu.Lock()
	if _, ok := c.items[key]; !ok {
		c.items[key] = &valueItem{
			Value:     value,
			startTime: time.Now(),
			Expires:   expire.Nanoseconds(),
		}
	}
	c.items[key].lastAccess = time.Now().UnixNano()
	c.mu.Unlock()
}

// Delete puts a value to a cache. If a key and value exists, overwrite it.
func (c *Cache_) Delete(key string) {
	c.mu.Lock()
	if _, ok := c.items[key]; !ok {
		delete(c.items, key)
	}
	c.mu.Unlock()
}
