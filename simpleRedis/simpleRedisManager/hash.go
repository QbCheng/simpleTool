package simpleRedisManager

import (
	"context"
)

func (rd *RedisMng) HSet(ctx context.Context, index int, key string, field string, value string) (int, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return 0, err
	}
	ret, err := conn.HSet(ctx, key, field, value).Result()
	return int(ret), err
}

func (rd *RedisMng) HGet(ctx context.Context, index int, key, field string) (string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return "", err
	}
	r, err := conn.HGet(ctx, key, field).Result()
	return r, err
}

func (rd *RedisMng) HGetAll(ctx context.Context, index int, key string) (map[string]string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return nil, err
	}
	r, err := conn.HGetAll(ctx, key).Result()
	return r, err
}
