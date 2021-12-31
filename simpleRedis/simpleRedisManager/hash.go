package simpleRedisManager

import (
	"context"
)

func (rd *RedisMng) HSet(ctx context.Context, index int, key string, field string, value string) (int64, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return 0, err
	}
	return conn.HSet(ctx, key, field, value).Result()
}

func (rd *RedisMng) HGet(ctx context.Context, index int, key, field string) (string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return "", err
	}
	return conn.HGet(ctx, key, field).Result()
}

func (rd *RedisMng) HGetAll(ctx context.Context, index int, key string) (map[string]string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return nil, err
	}
	return conn.HGetAll(ctx, key).Result()
}
