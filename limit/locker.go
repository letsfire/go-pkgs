package limit

import (
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/letsfire/redigo/v2"
)

// Locker 分布式锁
type Locker struct {
	prefix  string
	local   *sync.RWMutex
	client  *redigo.Client
	lockers map[string]*sync.Mutex
}

func (rl *Locker) Wrap(key string, ttl int64, fn func() error) error {
	e := rl.Lock(key, ttl)
	defer rl.Unlock(key)
	if e != nil {
		return e
	}
	return fn()
}

// Lock 锁定指定key, ttl为未正常解锁情况下的最多独占时长
func (rl *Locker) Lock(key string, ttl int64) error {
	key = rl.buildKey(key)
	rl.getLocal(key).Lock()
	_, err := rl.client.
		Execute(func(c redis.Conn) (interface{}, error) {
			var get bool  // get the lock
			var err error // redis op error
			var dts int64 // dead timestamp
			for get == false && err == nil {
				ts := time.Now().Unix()
				get, err = redis.Bool(c.Do("SETNX", key, ts+ttl+1))
				if err == nil && get == false { // failed
					dts, err = redis.Int64(c.Do("GET", key))
					if err == nil && dts <= ts { // expired
						dts, err = redis.Int64(c.Do("GETSET", key, ts+ttl+1))
						if err == nil && dts <= ts {
							break // acquired
						}
					}
					time.Sleep(time.Second) // very second
				}
			}
			return nil, err
		})
	return err
}

func (rl *Locker) Unlock(key string) {
	key = rl.buildKey(key)
	defer rl.getLocal(key).Unlock()
	_, _ = rl.client.
		Execute(func(c redis.Conn) (interface{}, error) {
			return c.Do("DEL", key)
		})
}

func (rl *Locker) IsLocked(key string) error {
	dts, err := rl.client.
		Int64(func(c redis.Conn) (interface{}, error) {
			return c.Do("GET", rl.buildKey(key))
		})
	if err == nil && dts > time.Now().Unix() {
		return fmt.Errorf("the key [%s] is locked", key)
	}
	return err
}

func (rl *Locker) getLocal(key string) *sync.Mutex {
	rl.local.RLock()
	if l, ok := rl.lockers[key]; ok {
		rl.local.RUnlock()
		return l
	}
	rl.local.RUnlock()
	return rl.newLocal(key)
}

func (rl *Locker) newLocal(key string) *sync.Mutex {
	rl.local.Lock()
	defer rl.local.Unlock()
	if l, ok := rl.lockers[key]; ok {
		return l
	}
	rl.lockers[key] = new(sync.Mutex)
	return rl.lockers[key]
}

func (rl *Locker) buildKey(key string) string {
	return fmt.Sprintf("%s.limit.locker.%s", rl.prefix, key)
}

func NewLocker(client *redigo.Client, prefix string) *Locker {
	return &Locker{
		local:  new(sync.RWMutex),
		client: client, prefix: prefix,
		lockers: make(map[string]*sync.Mutex),
	}
}
