package serverplugin

import (
	"hank.com/goetcd/discover/server"
	"testing"
	"time"
)

func TestStartEtcdPlugin1(t *testing.T) {
	r := EtcdV3RegisterPlugin{
		EtcdServers:    []string{"localhost:2379", "localhost:2381", "localhost:2383"},
		BasePath:       "node",
		ServiceAddress: "192.168.1.102", //自动获取IP
		ServiceMeta: &server.ServiceMeta{ //服务器信息
			IP:         "",
			Endpoint:   "192.168.1.102",
			Weight:     30,
			CreateTime: time.Now(),
		},
		UpdateInterval: 5,
	}

	//启动服务
	r.Start()
}

func TestStartEtcdPlugin2(t *testing.T) {
	r := EtcdV3RegisterPlugin{
		EtcdServers:    []string{"localhost:2379", "localhost:2381", "localhost:2383"},
		BasePath:       "node",
		ServiceAddress: "192.168.1.105", //自动获取IP
		ServiceMeta: &server.ServiceMeta{ //服务器信息
			IP:         "",
			Endpoint:   "192.168.1.105",
			Weight:     30,
			CreateTime: time.Now(),
		},
		UpdateInterval: 5,
	}

	//启动服务
	r.Start()
}

func TestStartEtcdPlugin3(t *testing.T) {
	r := EtcdV3RegisterPlugin{
		EtcdServers:    []string{"localhost:2379", "localhost:2381", "localhost:2383"},
		BasePath:       "node",
		ServiceAddress: "192.168.1.118", //自动获取IP
		ServiceMeta: &server.ServiceMeta{ //服务器信息
			IP:         "",
			Endpoint:   "192.168.1.118",
			Weight:     30,
			CreateTime: time.Now(),
		},
		UpdateInterval: 5,
	}

	//启动服务
	r.Start()
}
