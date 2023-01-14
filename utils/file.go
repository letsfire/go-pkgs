package utils

import (
	"os"
)

// IsDir 存在且为目录
func IsDir(dir string) bool {
	fi, err := os.Stat(dir)
	if err == nil {
		return fi.IsDir()
	}
	return false
}

// IsFile 存在且为文件
func IsFile(file string) bool {
	fi, err := os.Stat(file)
	if err == nil {
		return fi.IsDir() == false
	}
	return false
}

// FileExist 文件是否存在
func FileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
