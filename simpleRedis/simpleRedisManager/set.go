package simpleRedisManager

import "context"

func (rd *RedisMng) SAdd(ctx context.Context, index int, key string, members ...interface{}) (int64, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return 0, err
	}
	return conn.SAdd(ctx, key, members).Result()
}

func (rd *RedisMng) SMembers(ctx context.Context, index int, key string) ([]string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return nil, err
	}
	return conn.SMembers(ctx, key).Result()
}
