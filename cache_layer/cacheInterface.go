package cache_layer

type CacheInterface interface {
	SetCache(key string, value interface{})
	GetCache(key string) (interface{}, bool)
	DeleteFromCache(key string)
}
