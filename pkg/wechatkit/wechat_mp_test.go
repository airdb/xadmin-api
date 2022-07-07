package wechatkit

import (
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
	"github.com/silenceper/wechat/v2/miniprogram/config"
)

func TestWechatMiniProgram_CodeUnlimit(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("SETEX")

	pool := &redis.Pool{
		Dial:    func() (redis.Conn, error) { return conn, nil },
		MaxIdle: 10,
	}

	config := &config.Config{
		AppID:     os.Getenv("WECHAT_APP_ID"),
		AppSecret: os.Getenv("WECHAT_APP_SECRET"),
	}

	type args struct {
		page  string
		scene string
	}
	tests := []struct {
		name         string
		wx           *WechatMiniProgram
		args         args
		wantResponse bool
		wantErr      bool
	}{
		{``, NewWechatMiniProgram(config, pool), args{`pages/redirect/wxmpcode`, `id=1143039`}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := tt.wx.CodeUnlimit(tt.args.page, tt.args.scene)
			if (err != nil) != tt.wantErr {
				t.Errorf("WechatMiniProgram.CodeUnlimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResponse && len(gotResponse) == 0 {
				t.Errorf("WechatMiniProgram.CodeUnlimit() = %v, want not empty", len(gotResponse))
			}
			if !tt.wantResponse && len(gotResponse) > 0 {
				t.Errorf("WechatMiniProgram.CodeUnlimit() = %v, want empty", len(gotResponse))
			}
		})
	}
}
