package wechatkit

import (
	"github.com/gomodule/redigo/redis"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
)

func NewWechat(pool *redis.Pool) *wechat.Wechat {
	wx := wechat.NewWechat()

	var rp cache.Redis
	rp.SetConn(pool)

	wx.SetCache(&rp)

	return wx
}
