package main

import (
	"context"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"

	"hank.com/etcd-3.3.12-annotated/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:2381", "localhost:2383"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()

	cli.Put(context.Background(), "foob", "8888888")

	//用for循环去监视watch
	for {
		rch := cli.Watch(context.Background(), "foob")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				logs.Info("%+v", wresp)
			}
		}
	}
}
