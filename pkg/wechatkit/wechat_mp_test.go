package wechatkit

import (
	"reflect"
	"testing"
)

func TestWechatMiniProgram_CodeUnlimit(t *testing.T) {
	type args struct {
		page  string
		scene string
	}
	tests := []struct {
		name         string
		wx           *WechatMiniProgram
		args         args
		wantResponse []byte
		wantErr      bool
	}{
		{``, NewWechatMiniProgram(NewWechat()), args{`pages/article/detail/index`, `id=1143039`}, []byte{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := tt.wx.CodeUnlimit(tt.args.page, tt.args.scene)
			if (err != nil) != tt.wantErr {
				t.Errorf("WechatMiniProgram.CodeUnlimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("WechatMiniProgram.CodeUnlimit() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
