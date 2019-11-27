package main

import (
	"context"
	"hank.com/etcd-3.3.12-annotated/clientv3"
	"log"
	"os"
	"time"
)

//创建租约注册服务
type Service struct {
	client        *clientv3.Client
	groupName           string
	leaseid 	clientv3.LeaseID
	ttl int64

	stop    chan error
}

//NewService -
func NewService(groupName string,addr []string, ttl int64) (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})

	if err != nil{
		log.Fatal(err)
		return nil,err
	}

	ser := &Service{
		client: cli,
		groupName:groupName,
		ttl:ttl,
	}

	return ser, nil
}

//Start -启动服务
func (ser *Service) Start() error{
	leaseKeepAliveResponse,err := ser.keepAlive()

	if err != nil {
		log.Fatal(err)
		return  err
	}

	//监听租约
	go ser.ListenLeaseRespChan(leaseKeepAliveResponse)

	return nil
}

//ListenLeaseRespChan -
func (ser *Service)ListenLeaseRespChan(leaseKeepAliveResponse <-chan *clientv3.LeaseKeepAliveResponse){
	for{
		select {
		case ka,ok := <- leaseKeepAliveResponse:
			if !ok {
				log.Printf("keep alive channel closed\n")
				ser.RevokeLease()
			} else {
				log.Printf("Recv reply from service:%s,ttl:%d",ser.groupName,ka.TTL)
			}
		}
	}
}

//setLease -保持租约
func (ser *Service) keepAlive()(<-chan *clientv3.LeaseKeepAliveResponse,error){
	lease := clientv3.NewLease(ser.client)

	//设置租约时间
	leaseResp, err := lease.Grant(context.TODO(), ser.ttl)
	if err != nil {
		log.Fatal(err)
		return nil,err
	}

	ser.leaseid = leaseResp.ID
	return lease.KeepAlive(context.TODO(), leaseResp.ID)
}

//通过租约 注册服务
func (ser *Service) Register(nodeName, endpoint string) error {
	kv := clientv3.NewKV(ser.client)
	key := ser.groupName + "/"+ nodeName
	val := endpoint
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(ser.leaseid))
	return err
}

func (ser *Service)UnRegister(nodeName string)error{
	kv := clientv3.NewKV(ser.client)
	_,err := kv.Delete(context.TODO(),nodeName)
	return err
}

//撤销租约
func (ser *Service) RevokeLease() error {
	time.Sleep(2 * time.Second)
	_, err := ser.client.Revoke(context.TODO(), ser.leaseid)
	return err
}

func (ser *Service)Stop(){
	ser.stop <- nil
}

func main() {
	ser, err := NewService("node",[]string{"localhost:2379", "localhost:2381", "localhost:2383"}, 5)

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	ser.Register("hank1","localhost:2379")
	ser.Register("hank2","localhost:2381")
	ser.Register("hank3","localhost:2383")

	//启动服务
	ser.Start()
	select {}
}
