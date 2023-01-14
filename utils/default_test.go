package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDefault(t *testing.T) {
	var s string
	SetDefault(&s, "test1")
	assert.Equal(t, "test1", s)
	SetDefault(&s, "test2")
	assert.Equal(t, "test1", s)

	var i int
	SetDefault(&i, 1)
	assert.Equal(t, 1, i)
	SetDefault(&i, 2)
	assert.Equal(t, 1, i)

	var ss []string
	SetDefault(&ss, []string{"test1"})
	assert.Equal(t, []string{"test1"}, ss)
	SetDefault(&ss, []string{"test1", "test2"})
	assert.Equal(t, []string{"test1"}, ss)
}

func TestFillSlice(t *testing.T) {
	ss := make([]string, 0)
	FillSlice(&ss, 2, "111")
	assert.Len(t, ss, 2)
	assert.Equal(t, ss[1], "111")

	is := make([]int64, 0)
	FillSlice(&is, 2, int64(1))
	assert.Len(t, is, 2)
	assert.Equal(t, is[1], int64(1))
}
