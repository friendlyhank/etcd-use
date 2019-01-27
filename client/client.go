package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main(){
	cli,err := clientv3.New(clientv3.Config{
		Endpoints:[]string{"localhost:2379"},
		DialTimeout:5 * time.Second,
	})
	if err != nil{
		fmt.Println("connect failed,err",err)
		return
	}

	fmt.Println("connect success")

	defer cli.Close()

	//设置1秒超时，访问etcd有超时控制
	ctx,cancle:= context.WithTimeout(context.Background(),10 * time.Second)
	_,err = cli.Put(ctx,"foo","sample_value")
	//操作完毕取消
	cancle()

	if nil != err{
		fmt.Println("put failed,err",err)
		return
	}

	//取值,设置1秒超时
	ctx,cancle = context.WithTimeout(context.Background(),10 * time.Second)
	resp,err := cli.Get(ctx,"foo")

	if nil != err{
		fmt.Println("put failed,err",err)
		return
	}

	for _,ev := range resp.Kvs{
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}
