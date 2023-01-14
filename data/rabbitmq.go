package data

import (
	"fmt"
	"strings"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type Rabbitmq struct {
	Address  string // host:port
	Username string // 认证账号
	Password string // 安全密码

	locker sync.RWMutex
	vhosts map[string]*amqp091.Connection
}

func (r *Rabbitmq) GetConn(vhost string) *amqp091.Connection {
	r.locker.RLock()
	if v, ok := r.vhosts[vhost]; ok && !v.IsClosed() {
		r.locker.RUnlock()
		return v
	}
	r.locker.RUnlock()
	r.locker.Lock()
	defer r.locker.Unlock()
	return r.newConn(vhost)
}

func (r *Rabbitmq) Close() {
	r.locker.Lock()
	defer r.locker.Unlock()
	for vhost := range r.vhosts {
		r.vhosts[vhost].Close()
		delete(r.vhosts, vhost)
	}
}

func (r *Rabbitmq) newConn(vhost string) *amqp091.Connection {
	if v, ok := r.vhosts[vhost]; ok {
		v.Close()
	}
	url := fmt.Sprintf(
		"amqp://%s:%s@%s/%s",
		r.Username, r.Password, r.Address,
		strings.TrimLeft(vhost, "/"),
	)
	conn, err := amqp091.Dial(url)
	throwError(err) // 确保连接成功
	if r.vhosts == nil {
		r.vhosts = make(map[string]*amqp091.Connection)
	}
	r.vhosts[vhost] = conn
	return conn
}
