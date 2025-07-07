package xcache

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type CacheItem struct {
	Ctx context.Context

	Key   string
	Value interface{}

	// TTL is the cache expiration time.
	// Default TTL is 1 hour.
	TTL time.Duration

	// Do returns value to be cached.
	Do func(*CacheItem) (interface{}, error)

	// SetXX only sets the key if it already exists.
	SetXX bool

	// SetNX only sets the key if it does not already exist.
	SetNX bool

	// SkipLocalCache skips local cache as if it is not set.
	SkipLocalCache bool
}

type Cache interface {
	Delete(ctx context.Context, key string) error
	DeleteFromLocalCache(key string)
	Exists(ctx context.Context, key string) bool
	Get(ctx context.Context, key string, value interface{}) error
	GetSkippingLocalCache(ctx context.Context, key string, value interface{}) error
	Set(item *CacheItem) error
	Clean(ctx context.Context, pattern string) error
	KeysGet(ctx context.Context, pattern string) (map[string]interface{}, error)
	Clone(prefix string) Cache
	Keys(ctx context.Context, pattern string) ([]string, error)
}

type client struct {
	*cache.Cache
	redisc *redis.Client
	prefix string
}

type RedisConfig struct {
	Address string `envconfig:"REDIS_ADDRESS"`
	Prefix  string `envconfig:"REDIS_PREFIX" default:"/xcache/"`
}

func New(cfg *RedisConfig) Cache {
	redisc := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
	})

	pcache := cache.New(&cache.Options{
		Redis: redisc,
	})
	return &client{Cache: pcache, redisc: redisc, prefix: cfg.Prefix + "cache/"}
}

func (c *client) getKey(key string) string {
	return c.prefix + key
}

func (c *client) trimKey(key string) string {
	return strings.TrimPrefix(key, c.prefix)
}

func (c *client) Clone(prefix string) Cache {
	return &client{
		Cache:  c.Cache,
		redisc: c.redisc,
		prefix: prefix,
	}
}

func (c *client) Set(item *CacheItem) error {
	var do func(i *cache.Item) (interface{}, error)
	if item.Do != nil {
		do = func(i *cache.Item) (interface{}, error) {
			return item.Do(&CacheItem{
				Ctx: i.Ctx,
				Key: c.trimKey(i.Key),
				TTL: i.TTL})
		}
	}
	return c.Cache.Set(&cache.Item{
		Ctx:            item.Ctx,
		Key:            c.getKey(item.Key),
		Value:          item.Value,
		TTL:            item.TTL,
		Do:             do,
		SetXX:          item.SetXX,
		SetNX:          item.SetNX,
		SkipLocalCache: item.SkipLocalCache,
	})
}

func (c *client) Get(ctx context.Context, key string, value interface{}) error {
	return c.Cache.Get(ctx, c.getKey(key), value)
}

func (c *client) GetSkippingLocalCache(ctx context.Context, key string, value interface{}) error {
	return c.Cache.GetSkippingLocalCache(ctx, c.getKey(key), value)
}

func (c *client) Delete(ctx context.Context, key string) error {
	return c.Cache.Delete(ctx, c.getKey(key))
}

func (c *client) Exists(ctx context.Context, key string) bool {
	return c.Cache.Exists(ctx, c.getKey(key))
}

func (c *client) Clean(ctx context.Context, pattern string) error {
	res := c.redisc.Keys(ctx, c.getKey(pattern))
	if res.Err() != nil {
		return res.Err()
	}
	keys, err := res.Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		err := c.Cache.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) KeysGet(ctx context.Context, pattern string) (map[string]interface{}, error) {
	keyres := c.redisc.Keys(ctx, c.getKey(pattern))
	if keyres.Err() != nil {
		return nil, keyres.Err()
	}
	keys, err := keyres.Result()
	if err != nil {
		return nil, err
	}
	vres := c.redisc.MGet(ctx, keys...)
	if vres.Err() != nil {
		return nil, vres.Err()
	}
	res, err := c.redisc.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	var mres = make(map[string]interface{})
	for i, key := range keys {
		mres[key] = res[i]
	}
	return mres, nil
}

func (c *client) Keys(ctx context.Context, pattern string) ([]string, error) {
	keyres := c.redisc.Keys(ctx, c.getKey(pattern))
	if keyres.Err() != nil {
		return nil, keyres.Err()
	}
	return keyres.Result()
}
