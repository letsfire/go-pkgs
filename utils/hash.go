package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

var md5_ = md5.New()
var sha256_ = sha256.New()

// MD5 计算MD5值
func MD5(bytes []byte) string {
	return HashString(bytes, md5_)
}

// SHA256 计算SHA256值
func SHA256(bytes []byte) string {
	return HashString(bytes, sha256_)
}

func StringMD5(str string) string {
	return MD5([]byte(str))
}

func StringSHA256(str string) string {
	return SHA256([]byte(str))
}

// HashString 计算Hash字符串
func HashString(bytes []byte, hash hash.Hash) string {
	hash.Reset()
	hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}
