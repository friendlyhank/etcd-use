package discovery

import (
	"context"
	"fmt"
	"time"

	"hank.com/etcd-3.3.12-annotated/clientv3"
)

//创建租约注册服务
type ServiceReg struct {
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	canclefunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
}

func NewServiceReg(addr []string, timeNum int64) (*ServiceReg, error) {
	config := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}

	var (
		client *clientv3.Client
	)

	if clientTem, err := clientv3.New(config); err == nil {
		client = clientTem
	} else {
		return nil, err
	}

	ser := &ServiceReg{
		client: client,
	}

	if err := ser.setLease(timeNum);err != nil{
		return nil,err
	}

	return ser,nil
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
	ser.leaseResp  = leaseResp
	ser.keepAliveChan = leaseRespChan
	ser.canclefunc = cancelFunc

	return nil
}

//监听续租情况
func (ser *ServiceReg)ListenLeaseRespChan(){
	for{
		select{
			case leaseKeepResp := <- ser.keepAliveChan:
				if leaseKeepResp == nil{
					fmt.Printf("已经关闭续租功能\n")
					return
				}else{
					fmt.Printf("续租成功\n")
				}
		}
	}
}

//通过租约 注册服务
func (ser *ServiceReg)PutService(key,val string)error{
	kv := clientv3.NewKV(ser.client)
	_,err :=kv.Put(context.TODO(),key,val,clientv3.WithLease(ser.leaseResp.ID))
	return err
}

//撤销租约
func (ser *ServiceReg)RevokeLease()error{
	ser.canclefunc()
	time.Sleep(2 *time.Second)
	_,err := ser.lease.Revoke(context.TODO(),ser.leaseResp.ID)
	return err
}

func main(){
	ser,_ := NewServiceReg([]string{"127.0.0.1:2379"},5)
	//注册服务
	ser.PutService("/node/111","hello")
	select{}
}


