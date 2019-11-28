package main

import (
	"hank.com/goetcd/client/discover"
	"log"
	"os"
)

func main() {
	ser, err := discover.NewService([]string{"localhost:2379", "localhost:2381", "localhost:2383"}, 5)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	ser.PutService("node","hank1","localhost:2379")

	//启动服务
	ser.Start()
}
