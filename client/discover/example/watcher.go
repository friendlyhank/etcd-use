package main

import (
	"hank.com/goetcd/client/discover"
	"log"
	"time"
)

func main() {
	m, _ := discover.NewMaster("node/",[]string{"localhost:2379", "localhost:2381", "localhost:2383"},)
	for{
		m.Nodes.Range(func(k, v interface{})bool{
			log.Printf("node:%s, ip=%s endpoint=%s\n", k, v.(*discover.Node).ServiceMeta.IP, v.(*discover.Node).ServiceMeta.Endpoint)
			return true
		})
		time.Sleep(time.Second*5)
	}
}