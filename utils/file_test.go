package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	assert.True(t, FileExist("./file.go"))
	assert.False(t, FileExist("./file_not_exist.go"))
}
