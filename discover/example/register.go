package main

import (
	"errors"
	"hank.com/goetcd/discover"
	"hank.com/goetcd/discover/util"
	"log"
	"os"
	"time"
)

func main() {
	ser, err := discover.NewService([]string{"localhost:2379", "localhost:2381", "localhost:2383"}, 5)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	endpoint := util.GetGlobalUnicastIp()
	if endpoint == ""{
		err := errors.New("未获取对应endpoint")
		log.Fatal(err)
		os.Exit(2)
	}

	serviceMeta := &discover.ServiceMeta{IP:"183.6.58.77",Endpoint:endpoint,CreateTime:time.Now()}
	ser.PutService("node","hank1",serviceMeta)

	//启动服务
	ser.Start()
}
