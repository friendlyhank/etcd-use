package discover

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type Master struct {
	groupName string
	Nodes  *sync.Map
	client     *clientv3.Client
}

//Node- 监听节点信息
type Node struct{
	Key string
	Name string
	ServiceMeta *ServiceMeta
}

//ServiceMeta-
type ServiceMeta struct {
	IP  string
	Endpoint string
	Weight int32  //权重
	CreateTime time.Time //创建时间
}

//NewMaster-
func NewMaster(groupName string,addr []string,) (*Master, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})

	if err != nil{
		log.Fatal(err)
		return nil,err
	}

	m :=&Master{
		groupName:groupName,
		client:     cli,
		Nodes:new(sync.Map),
	}

	go m.watcher()

	return m,nil
}

//watcher-
func (m *Master) watcher() {
	rch := m.client.Watch(context.TODO(), m.groupName, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				info := GetServiceMeta(ev.Kv)
				m.AddNode(string(ev.Kv.Key), info)
			case mvccpb.DELETE:
				m.DelNode(string(ev.Kv.Key))
			}
		}
	}
}

//GetServiceMeta-
func GetServiceMeta(kv *mvccpb.KeyValue)*ServiceMeta{
	info := &ServiceMeta{}
	err := json.Unmarshal(kv.Value,info)
	if err != nil{
		log.Println(err)
	}
	return info
}

//GetServiceList-
func (m *Master) GetServiceList() ([]string, error) {
	var list []string
	m.Nodes.Range( func(k,v interface{})bool{
		node := v.(*Node)
		list = append(list,node.ServiceMeta.Endpoint)
		return true
	})
	return list,nil
}

//AddNode-新增节点
func (m *Master) AddNode(key string, info *ServiceMeta) {
	node := &Node{
		Key:key,
		ServiceMeta:info,
	}
	m.Nodes.Store(node.Key,node)
	log.Println("set data key :", key)
}

//DelNode- 删除节点
func (m *Master)DelNode(key string){
	m.Nodes.Delete(key)
	log.Println("set delete key :", key)
}

func (m *Master) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			info :=GetServiceMeta(resp.Kvs[i])
			addrs = append(addrs, info.Endpoint)
		}
	}

	return addrs
}


