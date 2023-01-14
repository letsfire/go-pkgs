package utils

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
)

// JwtClaims JWT请求体
type JwtClaims struct {
	jwt.StandardClaims
	Data Extend `json:"data"`
}

// JwtSigner JWT生成器
type JwtSigner struct {
	key *rsa.PrivateKey
}

func (s *JwtSigner) Encode(claims JwtClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(s.key)
}

func NewJwtSigner(keyFile string) *JwtSigner {
	var signer = new(JwtSigner)
	data, err := ioutil.ReadFile(keyFile)
	if err == nil {
		signer.key, err = jwt.ParseRSAPrivateKeyFromPEM(data)
	}
	if err != nil {
		panic(fmt.Errorf("parse rsa private key failed, error = %s", err))
	}
	return signer
}

// JwtParser JWT解析器
type JwtParser struct {
	key     *rsa.PublicKey
	keyFunc jwt.Keyfunc
}

func (p *JwtParser) Decode(s string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(s, &JwtClaims{}, p.keyFunc)
	return token.Claims.(*JwtClaims), err
}

func NewJwtParser(keyFile string, keyFunc func(*jwt.Token) error) *JwtParser {
	var parser = new(JwtParser)
	data, err := ioutil.ReadFile(keyFile)
	if err == nil {
		parser.key, err = jwt.ParseRSAPublicKeyFromPEM(data)
		parser.keyFunc = func(t *jwt.Token) (interface{}, error) {
			if keyFunc == nil {
				return parser.key, nil
			}
			return parser.key, keyFunc(t)
		}
	}
	if err != nil {
		panic(fmt.Errorf("parse rsa public key failed, error = %s", err))
	}
	return parser
}
