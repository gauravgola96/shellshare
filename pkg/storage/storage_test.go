package storage

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

func Setup() {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.SetConfigFile("../../config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	viper.AutomaticEnv()
}

func TestNewCache(t *testing.T) {
	Cache, _ := NewCache()
	key := "keytest_123"

	Cache.Put(key, "test data", 2*time.Second)
	time.Sleep(5 * time.Second)
	res, time, err := Cache.Get(key)
	if err != ErrNilCache && err != nil {
		return
	}
	t.Log(res, " ", time)
}

func TestStorageMongo(t *testing.T) {
	Setup()
	_ = Initialize()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.Mongo.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Log(err)
		return
	}

}
