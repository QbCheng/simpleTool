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

func (rd *RedisMng) RPush(ctx context.Context, index int, key string, value ...interface{}) (int64, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return 0, err
	}
	return conn.RPush(ctx, key, value...).Result()
}

func (rd *RedisMng) LTrim(ctx context.Context, index int, key string, start, end int) (string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return "", err
	}
	return conn.LTrim(ctx, key, int64(start), int64(end)).Result()
}

func (rd *RedisMng) LRange(ctx context.Context, index int, key string, start, end int) ([]string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return nil, err
	}
	return conn.LRange(ctx, key, int64(start), int64(end)).Result()
}
