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
	err := S.Mongo.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Log(err)
		return
	}
}

func TestRegisterUserData(t *testing.T) {
	Setup()
	_ = Initialize()
	err := RegisterUserData(context.TODO(), User{UserId: "ABC", SSHKeys: "SSH_123"})
	if err != nil {
		t.Log(err)
		return
	}

}

func TestUpdateUserLastLogin(t *testing.T) {
	Setup()
	_ = Initialize()
	err := UpdateUserLastLogin(context.TODO(), "ABC")
	if err != nil {
		t.Log(err)
		return
	}
}

func TestUpdateDownloadDetail(t *testing.T) {
	Setup()
	_ = Initialize()
	success := "success"
	err := UpdateDownloadDetail(context.TODO(), Download{
		SSHKeys:      "SSH_456",
		BytesWritten: 123456,
		Status:       Status(success),
	})
	if err != nil {
		t.Log(err)
		return
	}
}
func TestGetUsers(t *testing.T) {
	Setup()
	_ = Initialize()
	u, err := GetUsers(context.TODO(), -1)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(u)
}
