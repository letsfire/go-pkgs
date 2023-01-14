package data

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/go-pg/pg/v10"
)

type Postgres struct {
	Address  string // host:port
	Username string // 认证账号
	Password string // 安全密码

	locker sync.RWMutex
	dbs    map[string]*pg.DB
}

func (p *Postgres) GetConn(db string) *pg.DB {
	p.locker.RLock()
	if v, ok := p.dbs[db]; ok && p.test(v) == nil {
		p.locker.RUnlock()
		return v
	}
	p.locker.RUnlock()
	p.locker.Lock()
	defer p.locker.Unlock()
	return p.newConn(db)
}

func (p *Postgres) Close() {
	p.locker.Lock()
	defer p.locker.Unlock()
	for db := range p.dbs {
		p.dbs[db].Close()
		delete(p.dbs, db)
	}
}

func (p *Postgres) newConn(db string) *pg.DB {
	if v, ok := p.dbs[db]; ok {
		v.Close()
	}
	defaultTimeout := 10 * time.Second
	conn := pg.Connect(&pg.Options{
		Addr:         p.Address,
		User:         p.Username,
		Password:     p.Password,
		Database:     db,
		PoolSize:     0, // 暂时默认
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   defaultTimeout,
				KeepAlive: 15 * time.Second,
			}
			return netDialer.DialContext(ctx, network, addr)
		},
	})
	conn.AddQueryHook(&postgresLogger{slowLine: 5 * time.Second})
	throwError(p.test(conn)) // 确保连接成功
	if p.dbs == nil {
		p.dbs = make(map[string]*pg.DB)
	}
	p.dbs[db] = conn
	return conn
}

func (p *Postgres) test(db *pg.DB) error {
	_, err := db.Exec("SELECT 1")
	return err
}
