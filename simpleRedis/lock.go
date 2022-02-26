package simpleRedis

import (
	"context"
	"github.com/bsm/redislock"
	"time"
)

/*
ExponentialBackoffObtain
重试时间为2 * n毫秒(n表示次数).
可设置 最小值 和 最大值, 建议最小值不小于16ms.
*/
func (rm *RedisMng) ExponentialBackoffObtain(ctx context.Context, key string, ttl time.Duration, mateData string) (*redislock.Lock, error) {
	if !rm.RedisLockIsOpen() {
		return nil, ErrorDistributedNotOpen
	}
	option := &redislock.Options{
		RetryStrategy: redislock.ExponentialBackoff(20*time.Millisecond, 100*time.Millisecond),
		Metadata:      mateData,
	}
	return rm.distributed.Obtain(ctx, key, ttl, option)
}
