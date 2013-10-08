package cache

import (
    "github.com/darkhelmet/env"
)

type Cache interface {
    Get(key string) (string, error)
    Set(key string, data string, ttl int)
    Fetch(key string, ttl int, f func() string) string
}

var impl Cache = newDictCache()

func SetupMemcache(servers, username, password string) {
    impl = newMemcacheCache(servers, username, password)
}

func Get(key string) (string, error) {
    return impl.Get(key)
}

func Set(key, data string, ttl int) {
    impl.Set(key, data, ttl)
}

func Fetch(key string, ttl int, f func() string) string {
    return impl.Fetch(key, ttl, f)
}
