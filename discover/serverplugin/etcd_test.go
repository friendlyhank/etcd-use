package serverplugin

import (
	"hank.com/goetcd/discover/server"
	"hank.com/goetcd/discover/util"
	"testing"
	"time"
)

func TestStartEtcdPlugin(t *testing.T){
	r := EtcdV3RegisterPlugin{
		EtcdServers: []string{"localhost:2379", "localhost:2381", "localhost:2383"},
		BasePath:"node",
		ServiceAddress:util.GetGlobalUnicastIp(),//自动获取IP
		ServiceMeta:&server.ServiceMeta{//服务器信息
			IP:"",
			Endpoint:util.GetGlobalUnicastIp(),
			Weight:30,
			CreateTime:time.Now(),
		},
		UpdateInterval:20,
	}

	//启动服务
	r.Start()
}
