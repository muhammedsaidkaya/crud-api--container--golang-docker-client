package cache_layer

import (
	"github.com/docker/docker/api/types"
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache struct {
	inMemory *cache.Cache
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) CacheInterface {
	return &Cache{inMemory: cache.New(defaultExpiration, cleanupInterval)}
}

func (_cache Cache) SetCache(key string, value interface{}) {
	_cache.inMemory.Set(key, value, cache.DefaultExpiration)
}

func (_cache Cache) GetCache(key string) (interface{}, bool) {
	var container types.Container
	var found bool
	data, found := _cache.inMemory.Get(key)
	if found {
		container = data.(types.Container)
	}
	return container, found
}

func (_cache Cache) DeleteFromCache(key string) {
	_cache.inMemory.Delete(key)
}
