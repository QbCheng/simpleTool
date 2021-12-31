package simpleRedisManager

import "context"

func (rd *RedisMng) SAdd(ctx context.Context, index int, key string, members ...interface{}) (int, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return 0, err
	}
	r, err := conn.SAdd(ctx, key, members).Result()
	return int(r), err
}

func (rd *RedisMng) SMembers(ctx context.Context, index int, key string) ([]string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return nil, err
	}
	r, err := conn.SMembers(ctx, key).Result()
	return r, err
}
