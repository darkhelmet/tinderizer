package cache

import (
    "github.com/darkhelmet/env"
    "log"
    "os"
    "regexp"
)

type Cache interface {
    Get(key string) (string, error)
    Set(key string, data string, ttl int)
}

var impl Cache = newDictCache()
var logger *log.Logger

func SetupMemcache(servers, username, password string) {
    logger = log.New(os.Stdout, "[memcache] ", env.IntDefault("LOG_FLAGS", log.LstdFlags|log.Lmicroseconds))
    impl = newMemcacheCache(servers, username, password)
}

func SetupRedis(url, options string) {
    url = regexp.MustCompile(`^redis:`).ReplaceAllString(url, "tcp:")
    if options != "" {
        url += "?" + options
    }
    logger = log.New(os.Stdout, "[redis] ", env.IntDefault("LOG_FLAGS", log.LstdFlags|log.Lmicroseconds))
    impl = newRedisCache(url)
}

func Get(key string) (string, error) {
    return impl.Get(key)
}

func Set(key, data string, ttl int) {
    impl.Set(key, data, ttl)
}
