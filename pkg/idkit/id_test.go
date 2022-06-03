package idkit

import (
	"reflect"
	"testing"
	"time"
)

func TestNewWithTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		want Id
	}{
		{``, NewWithTime(now)},
	}
	for k, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if !reflect.DeepEqual(got.Machine(), tt.want.Machine()) {
				t.Errorf("New() = %v, want %v", got.Machine(), tt.want.Machine())
			}
			if !reflect.DeepEqual(got.Pid(), tt.want.Pid()) {
				t.Errorf("New() = %v, want %v", got.Pid(), tt.want.Pid())
			}
			if got.Counter()-tt.want.Counter() != int32(k)+1 {
				t.Errorf("New() = %v, want %v", got.Counter(), tt.want.Counter())
			}
			if got.Time().Unix()-now.Unix() > 1 {
				t.Errorf("New() = %v, want %v", got.Time().Unix(), now.Unix())
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkNewWithTime(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		NewWithTime(now)
	}
}
