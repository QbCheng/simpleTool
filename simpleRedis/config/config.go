package config

import "time"

type RedisConf struct {
	Addr     []string // Redis地址 ip:端口
	Password string   // Redis账号
	Db       int      // Redis库

	PoolSize        int           // Redis连接池大小
	MinIdleConnects int           // 最小空闲连接.
	IdleTimeout     time.Duration // 空闲链接超时时间

	MaxRetries      int           // 最大重试次数
	MinRetryBackoff time.Duration // 重试策略. 最短重连时间
	MaxRetryBackoff time.Duration // 重试策略. 最大重连时间

	DialTimeout  time.Duration // 连接超时时间
	ReadTimeout  time.Duration // 读超时
	WriteTimeout time.Duration // 写超时
}
