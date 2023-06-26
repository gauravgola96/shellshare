package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNilCache = errors.New("nil cache")
	S           *Storage
)

type Storage struct {
	Cache *Cache_
	Mongo *mongo.Client
}

// Cache_ is a struct for caching.
type Cache_ struct {
	items map[string]*ValueItem
	mu    sync.Mutex
}

type ValueItem struct {
	Message    string
	FileName   string
	StartTime  time.Time
	lastAccess int64
	Expires    int64
}

func Initialize() error {
	if S == nil {
		S = &Storage{}
	}
	m, err := InitializeMongo()
	if err != nil {
		return err
	}
	S.Mongo = m

	c, err := NewCache()
	if err != nil {
		return err
	}
	S.Cache = c
	return nil
}

func InitializeMongo() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", viper.GetString("storage.mongo.host"), viper.GetInt("storage.mongo.port"))))
	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	}
	//}()
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, err
}

func NewCache() (*Cache_, error) {
	subLogger := log.With().Str("module", "storage.cache").Logger()
	Cache := &Cache_{items: make(map[string]*ValueItem)}
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
	return Cache, nil
}

// Get gets a value from a cache.
func (c *Cache_) Get(key string) (*ValueItem, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.items[key]; ok {
		v.lastAccess = time.Now().UnixNano()
		return v, nil
	}
	return nil, ErrNilCache
}

// Put puts a value to a cache. If a key and value exists, overwrite it.
func (c *Cache_) Put(key string, value ValueItem, expire time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.items[key]; !ok {

		value.StartTime = time.Now()
		value.Expires = expire.Nanoseconds()
		c.items[key] = &value
	}
	c.items[key].lastAccess = time.Now().UnixNano()
}

// Delete Deletes a value to a cache.
func (c *Cache_) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.items[key]; ok {
		delete(c.items, key)
	}
}
