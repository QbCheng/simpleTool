package simpleRedisManager

import "context"

func (rd *RedisMng) LPop(ctx context.Context, index int, key string) (string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return "", err
	}
	r, err := conn.LPop(ctx, key).Result()
	return r, err
}

func (rd *RedisMng) RPush(ctx context.Context, index int, key string, value ...interface{}) (int, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return 0, err
	}
	r, err := conn.RPush(ctx, key, value...).Result()
	return int(r), err
}

func (rd *RedisMng) LTrim(ctx context.Context, index int, key string, start, end int) (string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return "", err
	}
	r, err := conn.LTrim(ctx, key, int64(start), int64(end)).Result()
	return r, err
}

func (rd *RedisMng) LRange(ctx context.Context, index int, key string, start, end int) ([]string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return nil, err
	}
	r, err := conn.LRange(ctx, key, int64(start), int64(end)).Result()
	return r, err
}
