package limit

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/letsfire/redigo/v2"
)

// Error 错误限制
type Error struct {
	prefix   string
	denyErr  error
	duration time.Duration
	client   *redigo.Client
}

func (e *Error) Run(max int, key string, fn func() error) error {
	if max <= 0 {
		return fn()
	}
	key = e.buildKey(key)
	val, err := e.client.
		Int(func(c redis.Conn) (interface{}, error) {
			return c.Do("GET", key)
		})
	if err == redis.ErrNil {
		val, err = 0, nil
	} else if err != nil {
		return err
	} else if val >= max {
		return e.denyErr
	}
	if err = fn(); err != nil {
		_, _ = e.client.
			Execute(func(c redis.Conn) (interface{}, error) {
				return c.Do("PSETEX", key, e.duration.Milliseconds(), val+1)
			})
		return err
	}
	return nil
}

func (e *Error) buildKey(key string) string {
	return fmt.Sprintf("%s.limit.error.%s", e.prefix, key)
}

func NewError(client *redigo.Client, prefix string, duration time.Duration, denyErr error) *Error {
	return &Error{client: client, prefix: prefix, duration: duration, denyErr: denyErr}
}
