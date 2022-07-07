package cachekit

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	onceRedisPool sync.Once

	redisPool *redis.Pool
)

func RedisPool(addr, passwd string, db int) *redis.Pool {
	onceRedisPool.Do(func() {
		redisPool = &redis.Pool{
			Dial: func() (redis.Conn, error) {
				dialOpts := []redis.DialOption{
					redis.DialDatabase(db),
					redis.DialPassword(passwd),
				}
				return redis.Dial("tcp", addr, dialOpts...)
			},
			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := conn.Do("PING")
				return err
			},
		}
	})

	return redisPool
}
