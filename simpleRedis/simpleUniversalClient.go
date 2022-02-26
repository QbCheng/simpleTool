package simpleRedis

import (
	"context"
	"errors"
	"github.com/QbCheng/simpleTool/simpleRedis/config"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

var (
	ErrorManagerInitializationFailed = errors.New(" Manager initialization failed ")
	ErrorDistributedNotOpen          = errors.New(" The distributed lock is not open ")
)

type Option func(*RedisMng)

type RedisMng struct {
	redis.UniversalClient // redis 统一客户端

	distributed *redislock.Client // redis 分布式锁
	option      option
}

func WithOpenRedisLockClient(openRedisLockClient bool) Option {
	return func(rm *RedisMng) {
		rm.option.openRedisLockClient = openRedisLockClient
	}
}

type option struct {
	openRedisLockClient bool
}

func NewRedisMng(conf *config.RedisConf, options ...Option) (*RedisMng, error) {
	m := new(RedisMng)
	var err error
	m.UniversalClient, err = create(conf)
	if err != nil {
		return nil, err
	}

	for i := range options {
		options[i](m)
	}

	// 初始化 redis 分布式锁
	m.initRedisLock()

	return m, nil
}

func (rm *RedisMng) RedisLockIsOpen() bool {
	return rm.option.openRedisLockClient && rm.distributed != nil
}

func (rm *RedisMng) initRedisLock() {
	if !rm.option.openRedisLockClient {
		return
	}
	rm.distributed = redislock.New(rm.UniversalClient)
}

func create(conf *config.RedisConf) (redis.UniversalClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.Addr,
		DB:       conf.Db,
		Password: conf.Password,

		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConnects,
		IdleTimeout:  conf.IdleTimeout,

		MaxRetries:      conf.MaxRetries,
		MinRetryBackoff: conf.MinRetryBackoff,
		MaxRetryBackoff: conf.MaxRetryBackoff,

		DialTimeout:  conf.DialTimeout,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, ErrorManagerInitializationFailed
	}
	return client, err
}

// KeyNil key 不存在.
func KeyNil(err error) bool {
	if err == redis.Nil {
		return true
	}
	return false
}
