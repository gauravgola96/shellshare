package storage

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	_ = InitializeCache()
	key := "keytest_123"

	Cache.Put(key, "test data", 2*time.Second)
	time.Sleep(5 * time.Second)
	res, time, err := Cache.Get(key)
	if err != NilCache() && err != nil {
		return
	}
	t.Log(res, " ", time)
}
