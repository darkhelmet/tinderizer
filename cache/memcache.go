package cache

import (
    "github.com/bmizerany/mc"
    "io"
    "syscall"
    "time"
)

type mcCache struct {
    conn     *mc.Conn
    server   string
    username string
    password string
}

func newMemcacheCache(server, username, password string) Cache {
    c := &mcCache{nil, server, username, password}
    c.connect()
    c.auth()
    return c
}

func (c *mcCache) connect() {
    logger.Printf("connecting")
    if cn, err := mc.Dial("tcp", c.server); err != nil {
        logger.Panicf("connecting failed: %s", err)
    } else {
        c.conn = cn
    }
}

func (c *mcCache) auth() {
    if c.username == "" {
        return
    }
    logger.Printf("authenticating")
    if err := c.conn.Auth(c.username, c.password); err != nil {
        logger.Panicf("authentication failed: %s", err)
    }
}

func (c *mcCache) handleError(action, key string, err error) (retry bool) {
    logger.Printf("error in %s for key %s: %s", action, key, err)
    switch err {
    case io.EOF, syscall.ECONNRESET:
        logger.Printf("trying to reconnect")
        // Lost connection? Try reconnecting
        time.Sleep(1 * time.Second)
        c.connect()
        // And of course we have to auth again
        fallthrough
    case mc.ErrAuthRequired:
        c.auth()
        return true
    case mc.ErrNotFound:
        // Cool story bro
    default:
        logger.Panicf("error in %s for key %s: %s", action, key, err)
    }
    return false
}

func (c *mcCache) Get(key string) (string, error) {
    return c.rget(key, 10)
}

func (c *mcCache) rget(key string, limit int) (string, error) {
    value, _, _, err := c.conn.Get(key)
    if err != nil {
        if c.handleError("get", key, err) && limit > 0 {
            return c.rget(key, limit-1)
        }
    }
    return value, err
}

func (c *mcCache) Set(key, data string, ttl int) {
    c.rset(key, data, ttl, 10)
}

func (c *mcCache) rset(key, data string, ttl, limit int) {
    // Don't worry about errors, live on the edge
    if err := c.conn.Set(key, data, 0, 0, ttl); err != nil {
        if c.handleError("set", key, err) && limit > 0 {
            c.rset(key, data, ttl, limit-1)
        }
    }
}
