package utils

import (
	"bytes"
	"math/rand"
	"strings"
)

const defaultSeed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomString 生成随机字符串
func RandomString(num int, seed string) string {
	return string(RandomBytes(num, seed))
}

// RandomBytes 生成随机字节集
func RandomBytes(num int, seed string) []byte {
	if seed == "" {
		seed = defaultSeed
	}
	buf := make([]byte, num)
	for i := 0; i < num; i++ {
		buf[i] = seed[rand.Intn(len(seed))]
	}
	return buf
}

// CamelToSnake 驼峰转蛇形
func CamelToSnake(s string) string {
	var flag bool
	buf := &bytes.Buffer{}
	buf.Grow(len(s) * 2)
	for i := 0; i < len(s); i++ {
		v := s[i]
		if flag && v >= 'A' && v <= 'Z' {
			buf.WriteByte('_')
		}
		if v != '_' {
			flag = true
		} else {
			flag = false
		}
		buf.WriteByte(v)
	}
	return strings.ToLower(buf.String())
}

// SnakeToCamel 蛇形转驼峰
func SnakeToCamel(s string) string {
	var flag = true
	buf := &bytes.Buffer{}
	buf.Grow(len(s))
	for i := 0; i < len(s); i++ {
		v := s[i]
		if v == '_' {
			flag = true
			continue
		} else if flag {
			v -= 32
			flag = false
		}
		buf.WriteByte(v)
	}
	return buf.String()
}
