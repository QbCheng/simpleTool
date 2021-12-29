package simpleRedisManager

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	ErrorManagerInitializationFailed = errors.New(" Manager initialization failed ")
	ErrorManagerIndexNotExist        = errors.New(" The Redis database does not exist ")
)

type RedisMng struct {
	ClientMap map[int]*redis.Client
}

func NewRedisMng() *RedisMng {
	m := new(RedisMng)
	m.ClientMap = make(map[int]*redis.Client)
	return m
}

type RedisConf struct {
	Index    int    // 数据库分组索引
	Addr     string // Redis地址 ip:端口
	Password string // Redis账号
	Db       int    // Redis库

	PoolSize        int           // Redis连接池大小
	MinIdleConnects int           // 最小空闲连接.
	IdleTimeout     time.Duration // 空闲链接超时时间

	MaxRetries      int           // 最大重试次数
	MinRetryBackoff time.Duration // 重试策略. 最短重连时间
	MaxRetryBackoff time.Duration // 重试策略. 最大重连时间

	DialTimeout  time.Duration // 连接超时时间
	ReadTimeout  time.Duration // 读超时
	WriteTimeout time.Duration // 写超时
}

func (rd *RedisMng) Init(configs []*RedisConf) error {
	for _, v := range configs {
		rd.ClientMap[v.Index] = redis.NewClient(&redis.Options{
			Addr:     v.Addr,
			DB:       v.Db,
			Password: v.Password,

			PoolSize:     v.PoolSize,
			MinIdleConns: v.MinIdleConnects,
			IdleTimeout:  v.IdleTimeout,

			MaxRetries:      v.MaxRetries,
			MinRetryBackoff: v.MinRetryBackoff,
			MaxRetryBackoff: v.MaxRetryBackoff,

			DialTimeout:  v.DialTimeout,
			ReadTimeout:  v.ReadTimeout,
			WriteTimeout: v.WriteTimeout,
		})
		_, err := rd.ClientMap[v.Index].Ping(context.Background()).Result()
		if err != nil {
			return ErrorManagerInitializationFailed
		}
	}
	return nil
}

func (rd *RedisMng) getConn(index int) (*redis.Client, error) {
	conn, ok := rd.ClientMap[index]
	if !ok {
		return nil, ErrorManagerIndexNotExist
	}
	return conn, nil
}
