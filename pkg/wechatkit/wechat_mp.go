package wechatkit

import (
	stdCtx "context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/gomodule/redigo/redis"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
)

type WechatMiniProgram struct {
	*miniprogram.MiniProgram
	redisPool *redis.Pool
	log       log.Fields
}

func NewWechatMiniProgram(config *config.Config, pool *redis.Pool, logger log.Logger) *WechatMiniProgram {
	return &WechatMiniProgram{
		MiniProgram: NewWechat(pool).GetMiniProgram(config),
		redisPool:   pool,
		log:         logger.WithField("kit", "wechat"),
	}
}

//
func (wx *WechatMiniProgram) Code2SessionContext(ctx stdCtx.Context, code string) (string, error) {
	res, err := wx.GetAuth().Code2SessionContext(ctx, code)
	if err != nil {
		return "", err
	}

	cacheKey := fmt.Sprintf("wxmp_oid:%s", res.OpenID)
	cacheContent, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	rpConn := wx.redisPool.Get()

	// 缓存30天
	if data, err := json.Marshal(cacheContent); err == nil {
		rpConn.Do("SETEX", cacheKey, 86400*30, string(data))
	}

	return res.OpenID, nil
}

// CodeUnlimit 生成小程序二维码
// Refer: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
func (wx *WechatMiniProgram) CodeUnlimit(page, scene string) (response []byte, err error) {
	width := 280
	cacheKey := fmt.Sprintf("wxmp_code:%s:%s:%d", page, scene, width)

	cacheContent := struct {
		Data []byte `json:"data,omitempty"`
	}{}
	rpConn := wx.redisPool.Get()
	if data, err := redis.Bytes(rpConn.Do("GET", cacheKey)); err == nil {
		err = json.Unmarshal(data, &cacheContent)
		if err == nil && len(cacheContent.Data) > 0 {
			wx.log.Debug(nil, "wxmp CodeUnlimit use cache")
			return cacheContent.Data, nil
		}
	}

	// 如果返回正常则进缓存
	codeParams := qrcode.QRCoder{
		Scene: scene,
		Page:  page,
		Width: width,
		EnvVersion: func() string {
			stage := os.Getenv(`ENV`)
			switch stage {
			case `test`: // 开发、测试环境
				return `develop`
			case `release`: // 正式环境
				fallthrough
			default:
				return `release`
			}
		}(),
	}
	if os.Getenv(`ENV`) == "test" {
		codeParams.CheckPath = false
	}
	qc := wx.GetQRCode()
	cacheContent.Data, err = qc.GetWXACodeUnlimit(codeParams)

	if err != nil {
		wx.log.Debug(nil, "query CodeUnlimit", codeParams, err)
		return []byte{}, err
	}

	// 缓存30天
	if data, err := json.Marshal(cacheContent); err == nil {
		rpConn.Do("SETEX", cacheKey, 86400*30, string(data))
	}

	return cacheContent.Data, err
}
