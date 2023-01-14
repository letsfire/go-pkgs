package conf

import (
	"math/rand"
	"os"
	"time"
)

type Root struct {
	Flag   string // 应用标识(中文)
	Mode   string // debug,test,release
	Port   string // HTTP运行端口
	Envs   Envs   // 环境变量
	Wechat Wechat // 微信配置
	Aliyun Aliyun // 阿里云配置
}

func (r *Root) prepare() *Root {
	rand.Seed(time.Now().Unix())
	for key, val := range r.Envs {
		throwError(os.Setenv(key, val))
	}
	return r
}
