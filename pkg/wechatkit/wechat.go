package wechatkit

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
)

var (
	onceRedisPool sync.Once

	redisPool *redis.Pool
)

func NewWechat() *wechat.Wechat {
	wx := wechat.NewWechat()

	var rp cache.Redis
	rp.SetConn(GetCacheRedisPool())

	wx.SetCache(&rp)
	return wx
}

func GetCacheRedisPool() *redis.Pool {
	onceRedisPool.Do(func() {
		database, err := strconv.Atoi(os.Getenv(`WXMP_REDIS_DB`))
		if err != nil {
			panic(err)
		}
		redisPool = &redis.Pool{
			Dial: func() (redis.Conn, error) {
				dialOpts := []redis.DialOption{
					redis.DialDatabase(database),
					redis.DialPassword(os.Getenv(`WX_REDIS_PASSWD`)),
				}
				return redis.Dial("tcp", os.Getenv(`WX_REDIS_ADDR`), dialOpts...)
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
