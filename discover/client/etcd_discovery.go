package client

import (
	"context"
	"log"
	"time"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type EtcdV3Discovery struct {
	BasePath string
	pairs []*KVPair
	client   *clientv3.Client
}

//NewMaster-
func NewEtcdV3Discovery(basePath string,etcdAddr []string) (*EtcdV3Discovery, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddr,
		DialTimeout: 5 * time.Second,
	})

	if err != nil{
		log.Fatal(err)
		return nil,err
	}

	m :=&EtcdV3Discovery{
		BasePath:basePath,
		client:     cli,
	}

	go m.watcher()

	return m,nil
}

//watcher-
func (m *EtcdV3Discovery) watcher() {
	rch := m.client.Watch(context.TODO(), m.BasePath, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				m.AddKVPair(string(ev.Kv.Key),string(ev.Kv.Value))
			case mvccpb.DELETE:

			}
		}
	}
}

//GetServiceList-
func (m *EtcdV3Discovery) GetServiceList() ([]string, error) {
	var list []string
	for _,pair :=range  m.pairs{
		list = append(list,pair.Value)
	}
	return list,nil
}

//AddNode-新增节点
func (m *EtcdV3Discovery) AddKVPair(key string, value string) {
	m.pairs = append(m.pairs,&KVPair{Key:key,Value:value})
	log.Println("set data key：%v,value：%v", key,value)
}

//DelNode- 删除节点
func (m *EtcdV3Discovery)DelKVPair(key string){
	for i,pair := range m.pairs{
		if pair.Key == key{
			m.pairs=append(m.pairs[:i],m.pairs[i+1:]...)
			log.Println("set delete key :", key)
		}
	}
}

func (m *EtcdV3Discovery) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			addrs = append(addrs, string(v))
		}
	}

	return addrs
}
