package client

import (
	"context"
	"log"
	"time"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

//EtcdV3Discovery-
type EtcdV3Discovery struct {
	BasePath string
	pairs []*KVPair
	client   *clientv3.Client
}

//NewEtcdV3Discovery-
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
		pairs:make([]*KVPair,0),
	}

	pairs,_ := m.GetKVPairList()
	if len(pairs) != 0{
		m.pairs = pairs
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
				m.DelKVPair(string(ev.Kv.Key))
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

func (m *EtcdV3Discovery) GetKVPairList() ([]*KVPair,error) {
	resp,err := m.client.Get(context.Background(),m.BasePath,clientv3.WithPrefix())
	if err != nil{
		return nil,err
	}

	pairs := make([]*KVPair,0)
	for i := range resp.Kvs {
		if string(resp.Kvs[i].Key) == "" || string(resp.Kvs[i].Value)==""{
			continue
		}
		pairs = append(pairs, &KVPair{Key:string(resp.Kvs[i].Key),Value:string(resp.Kvs[i].Value)})
	}

	return pairs,nil
}
