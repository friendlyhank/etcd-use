package main

import (
	"context"
	"fmt"
	"time"
	//可以自行换为官方etcd
	"hank.com/etcd-3.3.12-annotated/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:2381", "localhost:2383"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed,err", err)
		return
	}

	fmt.Println("connect success")

	defer cli.Close()

	//设置1秒超时，访问etcd有超时控制
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)

	//Put接口
	_, err = cli.Put(ctx, "foo", "sample_value")
	if err != nil {
		if err == context.Canceled {
			//II ctx is canceled by another routine
		} else if err == context.DeadlineExceeded {
			//II ctx is attached with a deadl ine and i t exceeded
		} else {
			//II bad cluster endpoints , which are not etcd servers
		}
	}
	//操作完毕取消
	cancle()

	if nil != err {
		fmt.Println("put failed,err", err)
		return
	}

	//取值,设置1秒超时
	ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)

	//Get接口
	resp, err := cli.Get(ctx, "foo")

	if nil != err {
		fmt.Println("put failed,err", err)
		return
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}
