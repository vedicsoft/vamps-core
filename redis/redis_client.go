package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/vedicsoft/vamps-core/commons"
)

var rg *Connection

type Connection struct {
	clusterClient  *redis.ClusterClient
	client         *redis.Client
	ClusterEnabled bool
	Addresses      []string
}

func init() {
	rg = &Connection{
		ClusterEnabled: false,
	}
	rg.client = redis.NewClient(&redis.Options{
		Addr:       commons.ServerConfigurations.RedisConfigs.Address,
		Password:   "",
		DB:         0,
		MaxRetries: 4,
	})
}

func SetValue(key string, value string, expiration int64) error {
	if rg.ClusterEnabled {
		return rg.clusterClient.Set(key, value, time.Duration(expiration)*time.Second).Err()
	} else {
		return rg.client.Set(key, value, time.Duration(expiration)*time.Second).Err()
	}
}

func GetValue(key string) (string, error) {
	if rg.ClusterEnabled {
		val, err := rg.clusterClient.Get(key).Result()
		if redis.Nil == err {
			return "", nil
		} else {
			return val, err
		}
	} else {
		val, err := rg.client.Get(key).Result()
		if redis.Nil == err {
			return "", nil
		} else {
			return val, err
		}
	}

}
