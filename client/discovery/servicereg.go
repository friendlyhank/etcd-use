package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"hank.com/etcd-3.3.12-annotated/clientv3"
)

//创建租约注册服务
type ServiceReg struct {
	name           string
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	canclefunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

//NewServiceReg -
func NewServiceReg(addr []string, timeNum int64) (*ServiceReg, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})

	if err != nil{
		log.Fatal(err)
		return nil,err
	}

	ser := &ServiceReg{
		client: cli,
	}

	if err := ser.setLease(timeNum); err != nil {
		return nil, err
	}

	//监听续租情况
	go ser.ListenLeaseRespChan()

	return ser, nil
}

//setLease -设置租约
func (ser *ServiceReg) setLease(timeNum int64) error {
	lease := clientv3.NewLease(ser.client)

	//设置租约时间
	leaseResp, err := lease.Grant(context.TODO(), timeNum)
	if err != nil {
		return err
	}

	//设置租约续期
	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)

	if err != nil {
		return err
	}

	ser.lease = lease
	ser.leaseResp = leaseResp
	ser.keepAliveChan = leaseRespChan
	ser.canclefunc = cancelFunc

	return nil
}

//监听续租情况
func (ser *ServiceReg) ListenLeaseRespChan() {
	for {
		select {
		case ka,ok := <-ser.keepAliveChan:
			if !ok {
				log.Printf("keep alive channel closed\n")
				ser.RevokeLease()
				return
			} else {
				log.Printf("Recv reply from service:%s,ttl:%d",ser.name,ka.TTL)
			}
		}
	}
}

//通过租约 注册服务
func (ser *ServiceReg) PutService(name, val string) error {
	kv := clientv3.NewKV(ser.client)
	_, err := kv.Put(context.TODO(), name, val, clientv3.WithLease(ser.leaseResp.ID))
	return err
}

//撤销租约
func (ser *ServiceReg) RevokeLease() error {
	ser.canclefunc()
	time.Sleep(2 * time.Second)
	_, err := ser.lease.Revoke(context.TODO(), ser.leaseResp.ID)
	return err
}

func main() {
	ser, err := NewServiceReg([]string{"localhost:2379", "localhost:2381", "localhost:2383"}, 5)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//注册服务
	ser.PutService("node/localhost:2379", "localhost:2379")
	ser.PutService("node/localhost:2381","localhost:2381")
	ser.PutService("node/localhost:2383","localhost:2383")
	select {}
}
