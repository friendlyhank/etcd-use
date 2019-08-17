package discovery

import (
	"context"
	"fmt"
	"hank.com/etcd-3.3.12-annotated/clientv3"
	"hank.com/etcd-3.3.12-annotated/mvcc/mvccpb"
	"sync"
	"time"
)

type ClientDis struct{
	client *clientv3.Client
	serverList map[string]string
	lock	sync.Mutex
}

func NewClientDis(addr []string)(*ClientDis,error){
	conf := clientv3.Config{
		Endpoints:addr,
		DialTimeout:5 *time.Second,
	}

	if client,err := clientv3.New(conf);err == nil{
		return &ClientDis{
			client:client,
		},nil
	}else{
		return nil,err
	}
}

func (cds *ClientDis)GetService(prefix string)([]string,error){
	resp,err := cds.client.Get(context.Background(),prefix,clientv3.WithPrefix())
	if err != nil{
		return nil,err
	}
	addrs := cds.extractAddrs(resp)

	go cds.watcher(prefix)

	return addrs,nil
}

func (cds *ClientDis)watcher(prefix string){
	rch := cds.client.Watch(context.TODO(),prefix,clientv3.WithPrefix())
	for wresp := range rch{
		for _,ev := range wresp.Events{
			switch ev.Type {
			case mvccpb.PUT:
				cds.SetServiceList(string(ev.Kv.Key),string(ev.Kv.Value))
			case mvccpb.DELETE:
				cds.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (cds *ClientDis)extractAddrs(resp *clientv3.GetResponse)[]string{
	addrs := make([]string,0)
	if resp == nil || resp.Kvs == nil{
		return addrs
	}

	for i := range resp.Kvs{
		if v := resp.Kvs[i].Value;v != nil{
			key := string(resp.Kvs[i].Key)
			value := string(resp.Kvs[i].Value)
			cds.SetServiceList(key,value)
			addrs = append(addrs,string(v))
		}
	}

	return addrs
}

func (cds *ClientDis)SetServiceList(key,val string){
	cds.lock.Lock()
	defer cds.lock.Unlock()
	cds.serverList[key] = string(val)
	fmt.Println("set data key :",key,"val:",val)
}

func (cds *ClientDis)DelServiceList(key string){
	cds.lock.Lock()
	defer cds.lock.Unlock()
	delete(cds.serverList,key)
	fmt.Println("del data key:", key)
}

func main(){
	cds,_ := NewClientDis([]string{"127.0.0.1:2379"})
	cds.GetService("/node")
	select{}
}
