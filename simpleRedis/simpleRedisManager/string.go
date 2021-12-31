package simpleRedisManager

import (
	"context"
	"time"
)

func (rd *RedisMng) Set(ctx context.Context, index int, key string, value interface{}, expiration time.Duration) (string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return "", err
	}
	return conn.Set(ctx, key, value, expiration).Result()
}

func (rd *RedisMng) Get(ctx context.Context, index int, key string) (string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return "", err
	}
	return conn.Get(ctx, key).Result()
}

func (rd *RedisMng) Incr(ctx context.Context, index int, key string) (int64, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return 0, err
	}
	return conn.Incr(ctx, key).Result()
}
