package utils

import (
	"fmt"
	"testing"
)

func TestUrlAdd(t *testing.T) {
	var url1 = "https://www.baidu.com"
	fmt.Println(UrlAdd(url1, map[string]string{"q": "1"}))

	var url2 = "https://www.baidu.com?q=1"
	fmt.Println(UrlAdd(url2, map[string]string{"w": "1"}))
}
