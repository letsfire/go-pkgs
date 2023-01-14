package limit

import (
	"errors"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/letsfire/redigo/v2/mode/alone"
	"github.com/stretchr/testify/assert"
)

var redisClient = alone.NewClient(
	alone.Addr("mlszp.com:6379"),
	alone.DialOpts(
		redis.DialPassword("Mangkaixin666!"),
	),
)

func TestError(t *testing.T) {
	var denyErr = errors.New("deny")
	var firstErr = errors.New("forbidden")
	var errLimiter = NewError(redisClient, "test", time.Second, denyErr)
	var err1 = errLimiter.Run(1, "scene1", func() error {
		return firstErr
	})
	var err2 = errLimiter.Run(1, "scene1", func() error {
		return firstErr
	})
	assert.Equal(t, err1, firstErr)
	assert.Equal(t, err2, denyErr)
	time.Sleep(time.Second)
	var err3 = errLimiter.Run(1, "scene1", func() error {
		return firstErr
	})
	assert.Equal(t, err3, firstErr)
}
