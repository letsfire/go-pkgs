package test

import (
	"context"
	"fmt"
	"github.com/letsfire/go-pkgs/data"
	"testing"
	"time"
)

var ctx context.Context
var zin data.Zincsearch

func init() {
	ctx = context.TODO()
	zin = data.Zincsearch{
		Address:  "http://www.mlszp.com:4080",
		Username: "zincsearch",
		Password: "vvtime@123456!",
	}
}

func TestZincsearch(t *testing.T) {
	_ = zin.IndexWithID(ctx, "space", "1", map[string]interface{}{
		"name": "章子宸 时光 集锦",
		"desc": "记录 脚步",
	})
	_ = zin.IndexWithID(ctx, "space", "2", map[string]interface{}{
		"name": "章子宸 证书 集锦",
		"desc": "记录 时光 脚步",
	})
	time.Sleep(time.Millisecond * 300)
	res, num, err := zin.Search(ctx, "space", "时光脚步", "name", "desc")
	fmt.Println(res, num, err)
}
