package main

import (
	"hank.com/goetcd/discover"
	"hank.com/goetcd/discover/util"
)

func main() {
	r := discover.EtcdV3RegisterPlugin{
		ServiceAddress:util.GetGlobalUnicastIp(),//自动获取IP
		EtcdServers: []string{"localhost:2379", "localhost:2381", "localhost:2383"},
		GroupName:"node",
		UpdateInterval:20,
	}

	//启动服务
	r.Start()
}
