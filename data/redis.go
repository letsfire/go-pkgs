package data

import (
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/letsfire/redigo/v2"
	"github.com/letsfire/redigo/v2/mode/alone"
)

const (
	RedisCache   = 0 // 缓存
	RedisEvent   = 1 // 事件
	RedisLocker  = 2 // 加锁
	RedisStorage = 3 // 存储
)

type Redis struct {
	Address  string // host:port
	Username string // 认证账号
	Password string // 安全密码

	locker  sync.RWMutex
	clients map[int]*redigo.Client
}

func (r *Redis) Cache() *redigo.Client {
	return r.getConn(RedisCache)
}

func (r *Redis) Event() *redigo.Client {
	return r.getConn(RedisEvent)
}

func (r *Redis) Locker() *redigo.Client {
	return r.getConn(RedisLocker)
}

func (r *Redis) Storage() *redigo.Client {
	return r.getConn(RedisStorage)
}

func (r *Redis) getConn(db int) *redigo.Client {
	r.locker.RLock()
	if v, ok := r.clients[db]; ok && r.test(v) == nil {
		r.locker.RUnlock()
		return v
	}
	r.locker.RUnlock()
	r.locker.Lock()
	defer r.locker.Unlock()
	return r.newConn(db)
}

func (r *Redis) Close() {
	r.locker.Lock()
	defer r.locker.Unlock()
	for db := range r.clients {
		r.clients[db].Close()
		delete(r.clients, db)
	}
}

func (r *Redis) newConn(db int) *redigo.Client {
	if v, ok := r.clients[db]; ok {
		v.Close()
	}
	cli := alone.NewClient(
		alone.Addr(r.Address),
		alone.DialOpts(
			redis.DialUsername(r.Username),
			redis.DialPassword(r.Password),
			redis.DialDatabase(db),
		),
		alone.PoolOpts(
			redigo.MaxActive(0), // 暂时不限制
		),
	)
	throwError(r.test(cli)) // 确保连接成功
	if r.clients == nil {
		r.clients = make(map[int]*redigo.Client)
	}
	r.clients[db] = cli
	return cli
}

func (r *Redis) test(cli *redigo.Client) error {
	_, err := cli.Execute(
		func(c redis.Conn) (interface{}, error) { return c.Do("PING") },
	)
	return err
}
