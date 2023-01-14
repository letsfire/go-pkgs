package utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	var data = JwtClaims{
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Unix()},
		Data:           Extend{"user_id": 1},
	}
	s, e1 := NewJwtSigner("./rsa_private").Encode(data)
	c, e2 := NewJwtParser("./rsa_public", nil).Decode(s)
	assert.Nil(t, e1)
	assert.Nil(t, e2)
	assert.Equal(t, c.Data.Int("user_id"), data.Data["user_id"])
}
