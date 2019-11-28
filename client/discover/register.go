package discover

import (
	"context"
	"encoding/json"
	"errors"
	"hank.com/etcd-3.3.12-annotated/clientv3"
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
	ser.Register(ser.node.key,ser.node.serviceMeta)

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
				log.Printf("Recv reply from service:%s,ttl:%d",ser.groupName,ka.TTL)
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
func (ser *Service)PutService(groupName,key string,endpoint string){
	ser.groupName = groupName
	ser.node = &Node{
		key:key,
		serviceMeta:&ServiceMeta{
			Endpoint:endpoint,
		},
	}
}

//通过租约 注册服务
func (ser *Service) Register(nodeName string, serviceMeta *ServiceMeta) error {
	kv := clientv3.NewKV(ser.client)
	key := ser.groupName + "/"+ nodeName
	val,_ := json.Marshal(serviceMeta)
	_, err := kv.Put(context.TODO(), key, string(val), clientv3.WithLease(ser.leaseid))
	return err
}

//UnRegister- 取消监听服务
func (ser *Service)UnRegister(nodeName string)error{
	kv := clientv3.NewKV(ser.client)
	key := ser.groupName+ "/"+ nodeName
	_,err := kv.Delete(context.TODO(),key)
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
