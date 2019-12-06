package serverplugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"hank.com/goetcd/discover/server"
	"log"
	"time"
)

/**
 *serverplugin 可插拔的插件模式
 */

//创建租约注册服务
type EtcdV3RegisterPlugin struct {
	//service address,for example,tcp@127.0.0.1:8972,tcp@127.0.0.1:8973
	ServiceAddress string
	//etcd addresss
	EtcdServers []string
	//Registered services
	Services       []string
	BasePath           string
	ServiceMeta *server.ServiceMeta

	client       *clientv3.Client
	leaseid 	clientv3.LeaseID
	UpdateInterval int64

	stop    chan error
}

//NewEtcdClentV3 -
func (ser *EtcdV3RegisterPlugin)NewEtcdClentV3() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   ser.EtcdServers,
		DialTimeout: 5 * time.Second,
	})

	if err != nil{
		log.Fatal(err)
		return nil,err
	}

	return cli, nil
}

//Start -启动服务
func (ser *EtcdV3RegisterPlugin) Start() error{

	//NewEtcdV3
	c,err := ser.NewEtcdClentV3()
	if err != nil{
		log.Fatalf("cannot create etcd registry: %v",err)
		return err
	}
	ser.client = c

	//保持租约
	leaseKeepAliveResponse,err := ser.keepAlive()
	if err != nil {
		log.Fatal(err)
		return  err
	}

	//注册监听
	meta, err := json.Marshal(&ser.ServiceMeta)
	ser.Register(ser.GetNodePath(), string(meta))

	for{
		select {
		case err := <- ser.stop:
			ser.RevokeLease()
			return err
		case <-ser.client.Ctx().Done():
			return errors.New("server closed")
		case ka,ok := <- leaseKeepAliveResponse://租约
			if !ok {
				log.Printf("keep alive channel closed\n")
				ser.RevokeLease()
			} else {
				log.Printf("Recv reply from service:%s,ttl:%d",ser.ServiceAddress,ka.TTL)
			}
		}
	}

	return nil
}

//setLease -保持租约
func (ser *EtcdV3RegisterPlugin) keepAlive()(<-chan *clientv3.LeaseKeepAliveResponse,error){
	lease := clientv3.NewLease(ser.client)

	//设置租约时间
	leaseResp, err := lease.Grant(context.TODO(), ser.UpdateInterval)
	if err != nil {
		log.Fatal(err)
		return nil,err
	}

	ser.leaseid = leaseResp.ID
	return lease.KeepAlive(context.TODO(), leaseResp.ID)
}

//通过租约 注册服务
func (ser *EtcdV3RegisterPlugin) Register(key string,val string) error {
	kv := clientv3.NewKV(ser.client)

	_, err := kv.Put(context.TODO(), key,val, clientv3.WithLease(ser.leaseid))
	return err
}

//GetNodePath-
func (ser *EtcdV3RegisterPlugin)GetNodePath()string{
	return fmt.Sprintf("%s/%s",ser.BasePath,ser.ServiceAddress)
}

//UnRegister- 取消监听服务
func (ser *EtcdV3RegisterPlugin)UnRegister()error{
	kv := clientv3.NewKV(ser.client)

	nodePath := ser.GetNodePath()
	_,err := kv.Delete(context.TODO(),nodePath)
	return err
}

//撤销租约
func (ser *EtcdV3RegisterPlugin) RevokeLease() error {
	_, err := ser.client.Revoke(context.TODO(), ser.leaseid)
	return err
}

//Stop-
func (ser *EtcdV3RegisterPlugin)Stop(){
	ser.stop <- nil
}

