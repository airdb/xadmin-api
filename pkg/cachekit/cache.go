package cachekit

type Cache struct {
	config map[string]CacheConfig
}

type CacheConfig struct {
	Type     string
	Host     string
	Password string
}

func NewCache(config map[string]CacheConfig) *Cache {
	return &Cache{
		config: config,
	}
}

func (c Cache) Redis(db int) *Redis {
	for _, v := range c.config {
		if v.Type == "redis" {
			return NewRedis(v.Host, v.Password, db)
		}
	}

	return nil
}
