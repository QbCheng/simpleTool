package simpleRedisLock

import (
	"context"
	"errors"
	"github.com/bsm/redislock"
	"simpleTool/simpleRedis/simpleRedisManager"
	"sync"
	"time"
)

var (
	ErrorDefaultSimpleRedisLockNotInit = errors.New(" Error DefaultSimpleRedisLock Not Init")
)

var (
	defaultSimpleRedisLock *SimpleRedisLock
	once                   sync.Once
)

type SimpleRedisLock struct {
	*redislock.Client
}

func NewSimpleRedisLock(conf simpleRedisManager.RedisConf) (*SimpleRedisLock, error) {
	ret := &SimpleRedisLock{}
	client, err := simpleRedisManager.CreateRedisClient(&conf)
	if err != nil {
		return nil, simpleRedisManager.ErrorManagerInitializationFailed
	}
	ret.Client = redislock.New(client)
	return ret, nil
}

func DefaultInit(conf simpleRedisManager.RedisConf) error {
	var err error
	if defaultSimpleRedisLock == nil {
		once.Do(func() {
			defaultSimpleRedisLock, err = NewSimpleRedisLock(conf)
		})
	}
	return err
}

func ExponentialBackoffObtain(ctx context.Context, key string, ttl time.Duration, mateData string) (*redislock.Lock, error) {
	if defaultSimpleRedisLock == nil {
		return nil, ErrorDefaultSimpleRedisLockNotInit
	}
	option := &redislock.Options{
		RetryStrategy: redislock.ExponentialBackoff(20*time.Millisecond, 100*time.Millisecond),
		Metadata:      mateData,
	}
	return defaultSimpleRedisLock.Obtain(ctx, key, ttl, option)
}
