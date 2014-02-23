package cache

import (
    "github.com/xuyu/goredis"
)

type redisCache struct {
    redis *goredis.Redis
    url   string
}

func newRedisCache(url string) Cache {
    c := &redisCache{nil, url}
    c.connect()
    return c
}

func (c *redisCache) connect() {
    logger.Printf("connecting")
    redis, err := goredis.DialURL(c.url)
    if err != nil {
        logger.Panicf("connecting failed: %s", err)
    }
    c.redis = redis
}

func (c *redisCache) handleError(action, key string, err error) {
    if err != nil {
        logger.Panicf("error in %s for key %s: %s", action, key, err)
    }
}

func (c *redisCache) Get(key string) (string, error) {
    data, err := c.redis.Get(key)
    c.handleError("get", key, err)
    return string(data), nil
}

func (c *redisCache) Set(key, data string, ttl int) {
    err := c.redis.Set(key, data, ttl, 0, false, false)
    c.handleError("set", key, err)
}
