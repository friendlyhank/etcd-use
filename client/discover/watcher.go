package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"hank.com/etcd-3.3.12-annotated/clientv3"
	"hank.com/etcd-3.3.12-annotated/mvcc/mvccpb"
)

type Master struct {
	groupName string
	Nodes  *sync.Map
	client     *clientv3.Client
}

//Node- 监听节点信息
type Node struct{
	nodeName string
	endpoint string
}

//NewClientDis-
func NewClientDis(addr []string,groupName string) (*Master, error) {
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
				m.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				m.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

//GetServiceList-
func (m *Master) GetServiceList() ([]string, error) {
	resp, err := m.client.Get(context.Background(), m.groupName, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	addrs := m.extractAddrs(resp)

	return addrs, nil
}

//SetServiceList-
func (m *Master) SetServiceList(key, val string) {
	node := &Node{
		nodeName:key,
		endpoint:val,
	}
	m.Nodes.Store(node.nodeName,node)
	fmt.Println("set data key :", key, "val:", val)
}

func (m *Master) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			key := string(resp.Kvs[i].Key)
			value := string(resp.Kvs[i].Value)
			m.SetServiceList(key, value)
			addrs = append(addrs, string(v))
		}
	}

	return addrs
}

func (m *Master) DelServiceList(key string) {
	m.Nodes.Delete(key)
	fmt.Println("del data key:", key)
}

func main() {
	cds, _ := NewClientDis([]string{"localhost:2379", "localhost:2381", "localhost:2383"},"node/")
	serviceList,_ := cds.GetServiceList()
	fmt.Println(serviceList)
	select {}
}
