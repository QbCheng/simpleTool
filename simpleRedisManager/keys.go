package simpleRedisManager

import "context"

func (rd *RedisMng) Keys(ctx context.Context, index int, pattern string) ([]string, error) {
	conn, err := rd.getConn(index)
	if err != nil {
		return nil, err
	}
	r, err := conn.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return r, err
}
