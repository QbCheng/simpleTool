package simpleUniversalClient

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"simpleTool/simpleRedis/config"
	"strconv"
	"sync"
	"testing"
	"time"
)

var TestConfig = config.RedisConf{
	Addr:     []string{"192.168.110.230:6379"}, // Redis地址 ip:端口
	Password: "pass230",                        // Redis账号
	Db:       0,                                // Redis库

	PoolSize:        10,               // Redis连接池大小
	MinIdleConnects: 5,                // 最小空闲连接.
	IdleTimeout:     time.Minute * 10, // 空闲链接超时时间

	MaxRetries:      10,                     // 最大重试次数
	MinRetryBackoff: time.Millisecond * 16,  // 重试策略. 最短重连时间
	MaxRetryBackoff: time.Millisecond * 512, // 重试策略. 最大重连时间

	DialTimeout:  time.Second * 10, // 连接超时时间
	ReadTimeout:  time.Second * 30, // 读超时
	WriteTimeout: time.Second * 30, // 写超时
}

func TestSingle(t *testing.T) {
	commonClient, err := NewRedisMng(&TestConfig, WithOpenRedisLockClient(false))
	assert.NoError(t, err)

	lockClient, err := NewRedisMng(&TestConfig, WithOpenRedisLockClient(true))
	assert.NoError(t, err)

	var firstData string
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			locker, err := lockClient.ExponentialBackoffObtain(context.Background(), " my-lock", 30*time.Second, "wo is mateData")
			assert.NoError(t, err)
			defer func() {
				err := locker.Release(context.Background())
				assert.NoError(t, err)
			}()
			fmt.Println(locker.Metadata())
			if firstData == "" {
				firstData = strconv.Itoa(rand.Int())
				commonClient.Set(context.Background(), "firstData", firstData, 0)
			} else {
				value, err := commonClient.Get(context.Background(), "firstData").Result()
				assert.NoError(t, err)
				assert.Equal(t, value, firstData)
			}
		}()
	}
	wg.Wait()
}

var lockKeyFormat = "my-lock:%v"

func TestMulti(t *testing.T) {
	commonClient, err := NewRedisMng(&TestConfig, WithOpenRedisLockClient(false))
	assert.NoError(t, err)

	lockClient, err := NewRedisMng(&TestConfig, WithOpenRedisLockClient(true))
	assert.NoError(t, err)

	var firstData1 string
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			lock := fmt.Sprintf(lockKeyFormat, 1)
			locker, err := lockClient.ExponentialBackoffObtain(context.Background(), lock, 30*time.Second, "wo is mateData")
			assert.NoError(t, err)
			defer func() {
				err := locker.Release(context.Background())
				assert.NoError(t, err)
			}()
			fmt.Println(locker.Metadata())
			if firstData1 == "" {
				firstData1 = strconv.Itoa(rand.Int())
				commonClient.Set(context.Background(), "firstData1", firstData1, 0)
			} else {
				value, err := commonClient.Get(context.Background(), "firstData1").Result()
				assert.NoError(t, err)
				assert.Equal(t, value, firstData1)
			}
		}()
	}

	var firstData2 string
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			lock := fmt.Sprintf(lockKeyFormat, 2)
			locker, err := lockClient.ExponentialBackoffObtain(context.Background(), lock, 30*time.Second, "wo is mateData")
			assert.NoError(t, err)
			defer func() {
				err := locker.Release(context.Background())
				assert.NoError(t, err)
			}()
			fmt.Println(locker.Metadata())
			if firstData2 == "" {
				firstData2 = strconv.Itoa(rand.Int())
				commonClient.Set(context.Background(), "firstData2", firstData2, 0)
			} else {
				value, err := commonClient.Get(context.Background(), "firstData2").Result()
				assert.NoError(t, err)
				assert.Equal(t, value, firstData2)
			}
		}()
	}
	wg.Wait()
}
