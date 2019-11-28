package main

import (
	"hank.com/goetcd/client/discover"
	"log"
	"os"
	"time"
)

func main() {
	ser, err := discover.NewService([]string{"localhost:2379", "localhost:2381", "localhost:2383"}, 20)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	serviceMeta := &discover.ServiceMeta{IP:"183.6.58.77",Endpoint:"localhost:2379",CreateTime:time.Now()}
	ser.PutService("node","hank1",serviceMeta)

	//启动服务
	ser.Start()
}
