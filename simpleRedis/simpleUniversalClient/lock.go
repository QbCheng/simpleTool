package simpleUniversalClient

import (
	"context"
	"github.com/bsm/redislock"
	"time"
)

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
