package discover

import (
	"context"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

//创建租约注册服务
type Service struct {
	client        *clientv3.Client
	groupName           string
	leaseid 	clientv3.LeaseID
	ttl int64

	node *Node
	stop    chan error
}

//NewService -
func NewService(addr []string, ttl int64) (*Service, error) {
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
		ttl:ttl,
	}

	return ser, nil
}

//Start -启动服务
func (ser *Service) Start() error{

	//保持租约
	leaseKeepAliveResponse,err := ser.keepAlive()
	if err != nil {
		log.Fatal(err)
		return  err
	}

	//注册服务
	if ser.node == nil{
		err = errors.New("未找到设置节点")
		log.Fatal(err)
		return err
	}

	ser.Register()

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
					log.Printf("Recv reply from service:%s,ttl:%d",ser.node.Key,ka.TTL)
			}
		}
	}

	return nil
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

//PutService-
func (ser *Service)PutService(groupName,nodeName string,serviceMeta *ServiceMeta){
	ser.groupName = groupName
	ser.node = &Node{
		Key:ser.groupName + "/"+ nodeName,
		Name:nodeName,
		ServiceMeta:serviceMeta,
	}
}

//通过租约 注册服务
func (ser *Service) Register() error {
	kv := clientv3.NewKV(ser.client)
	val,_ := json.Marshal(ser.node.ServiceMeta)
	_, err := kv.Put(context.TODO(), ser.node.Key, string(val), clientv3.WithLease(ser.leaseid))
	return err
}

//UnRegister- 取消监听服务
func (ser *Service)UnRegister()error{
	kv := clientv3.NewKV(ser.client)
	_,err := kv.Delete(context.TODO(),ser.node.Key)
	return err
}

//撤销租约
func (ser *Service) RevokeLease() error {
	_, err := ser.client.Revoke(context.TODO(), ser.leaseid)
	return err
}

func (ser *Service)Stop(){
	ser.stop <- nil
}
