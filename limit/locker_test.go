package limit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocker(t *testing.T) {
	locker := NewLocker(redisClient, "test")
	assert.Nil(t, locker.Wrap("scene1", 1, func() error { return nil }))
	assert.Nil(t, locker.Wrap("scene1", 1, func() error { return nil }))
	assert.NotNil(t, locker.Wrap("scene2", 1, func() error { return locker.IsLocked("scene2") }))
}
