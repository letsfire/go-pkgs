package conf

import (
	"context"
	"sync"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	mpCfg "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	oaCfg "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/openplatform"
	opCfg "github.com/silenceper/wechat/v2/openplatform/config"
)

type Wechat struct {
	sync.Mutex
	MPConfig  mpCfg.Config // 小程序配置
	OAConfig  oaCfg.Config // 公众号配置
	OPConfig  opCfg.Config // 开放平台配置
	RedisHost string       // Redis地址
	RedisPass string       // Redis密码
	sdk       *wechat.Wechat
}

func (w *Wechat) getSDK() *wechat.Wechat {
	w.Lock()
	defer w.Unlock()
	if w.sdk != nil {
		return w.sdk
	}
	var useCache cache.Cache
	if w.RedisHost != "" {
		useCache = cache.NewRedis(context.TODO(), &cache.RedisOpts{
			Host: w.RedisHost, Password: w.RedisPass,
			Database: 0, MaxIdle: 5, IdleTimeout: 30,
		})
	} else {
		useCache = cache.NewMemory()
	}
	w.sdk = wechat.NewWechat()
	w.sdk.SetCache(useCache)
	return w.sdk
}

func (w *Wechat) MiniProgram() *miniprogram.MiniProgram {
	return w.getSDK().GetMiniProgram(&w.MPConfig)
}

func (w *Wechat) OfficialAccount() *officialaccount.OfficialAccount {
	return w.getSDK().GetOfficialAccount(&w.OAConfig)
}

func (w *Wechat) OpenPlatform() *openplatform.OpenPlatform {
	return w.getSDK().GetOpenPlatform(&w.OPConfig)
}
