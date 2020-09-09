package client

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"github.com/friendlyhank/etcd-use/discover/server"
	"log"
	"time"
)



//EtcdV3Discovery-
type EtcdV3Discovery struct {
	BasePath string
	nodes []*server.Node
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
		nodes:make([]*server.Node,0),
	}

	nodes,_ := m.GetReadyNodes()
	if len(nodes) != 0{
		m.nodes = nodes
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
				serviceMeta := &server.ServiceMeta{}
				_ = json.Unmarshal(ev.Kv.Value, serviceMeta)
				m.AddNode(string(ev.Kv.Key),serviceMeta)
			case mvccpb.DELETE:
				m.DelNode(string(ev.Kv.Key))
			}
		}
	}
}

//GetServiceList-
func (m *EtcdV3Discovery) GetServiceList() ([]string, error) {
	var list []string
	for _,node :=range  m.nodes{
		list = append(list,node.ServiceMeta.Endpoint)
	}
	return list,nil
}

//AddNode-新增节点
func (m *EtcdV3Discovery) AddNode(key string, serviceMeta *server.ServiceMeta) {
	m.nodes = append(m.nodes,&server.Node{Key: key,ServiceMeta:serviceMeta})
	log.Printf("set data key：%v,value：%+v", key,serviceMeta)
}

//DelNode- 删除节点
func (m *EtcdV3Discovery)DelNode(key string){
	for i,node := range m.nodes{
		if node.Key == key{
			m.nodes=append(m.nodes[:i],m.nodes[i+1:]...)
			log.Printf("delete key :%v", key)
		}
	}
}

func (m *EtcdV3Discovery) GetReadyNodes() ([]*server.Node,error) {
	resp,err := m.client.Get(context.Background(),m.BasePath,clientv3.WithPrefix())
	if err != nil{
		return nil,err
	}

	nodes := make([]*server.Node,0)
	for i := range resp.Kvs {
		if string(resp.Kvs[i].Key) == "" || string(resp.Kvs[i].Value)==""{
			continue
		}
		serviceMeta := &server.ServiceMeta{}
		_ = json.Unmarshal(resp.Kvs[i].Value, serviceMeta)
		nodes = append(nodes, &server.Node{Key: string(resp.Kvs[i].Key),ServiceMeta:serviceMeta})
	}

	return nodes,nil
}
