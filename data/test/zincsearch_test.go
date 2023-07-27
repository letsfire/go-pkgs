package test

import (
	"context"
	"github.com/letsfire/go-pkgs/data"
	"testing"
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
		"name": "章子宸证书集锦",
		"desc": "记录时光的脚步",
	})
	_ = zin.IndexWithID(ctx, "story", "1", map[string]interface{}{
		"name": "章子宸证书集锦",
		"desc": "记录时光的脚步",
		"addr": "就是一大段内容怎",
	})
	zin.Search(ctx, "space", "证书")
}
