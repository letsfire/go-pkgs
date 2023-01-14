package conf

import (
	"github.com/aliyun-sdk/sms-go"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Aliyun struct {
	RegionId  string
	AccessKey string
	SecretKey string
	UseVpcNet bool
}

func (a *Aliyun) SLS() sls.ClientInterface {
	endpoint := a.RegionId
	if a.UseVpcNet {
		endpoint += "-intranet"
	}
	return sls.CreateNormalInterface(endpoint+".log.aliyuncs.com", a.AccessKey, a.SecretKey, "")
}

func (a *Aliyun) STS() *sts.Client {
	cli, err := sts.NewClientWithAccessKey(a.RegionId, a.AccessKey, a.SecretKey)
	throwError(err)
	return cli
}

func (a *Aliyun) SMS() *sms.Client {
	cli, err := sms.New(a.AccessKey, a.SecretKey)
	throwError(err)
	return cli
}

func (a *Aliyun) OSS() *oss.Client {
	endpoint := "oss-" + a.RegionId
	if a.UseVpcNet {
		endpoint += "-internal"
	}
	cli, err := oss.New(endpoint+".aliyuncs.com", a.AccessKey, a.SecretKey)
	throwError(err)
	return cli
}
