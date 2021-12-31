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
	r, err := conn.Set(ctx, key, value, expiration).Result()
	return r, err
}

func (rd *RedisMng) Get(ctx context.Context, index int, key string) (string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return "", err
	}
	r, err := conn.Get(ctx, key).Result()
	return r, err
}

func (rd *RedisMng) Incr(ctx context.Context, index int, key string) (int, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return 0, err
	}
	r, err := conn.Incr(ctx, key).Result()
	return int(r), err
}
