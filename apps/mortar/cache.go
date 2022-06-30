package mortar

import (
	"github.com/airdb/xadmin-api/pkg/cachekit"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"go.uber.org/fx"
)

// CacheFxOption registers Cache
func CacheFxOption() fx.Option {
	return fx.Options(
		fx.Provide(NewCache),
	)
}

func NewCache(config cfg.Config) *cachekit.Cache {
	cfgs := make(map[string]cachekit.CacheConfig)
	err := config.Get("xadmin.cache").Unmarshal(&cfgs)
	if err != nil {
		panic(err)
	}

	return cachekit.NewCache(cfgs)
}
