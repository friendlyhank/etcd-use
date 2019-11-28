package main

import(
	"hank.com/goetcd/client/discover"
	"fmt"
)

func main() {
	cds, _ := discover.NewClientDis("node/",[]string{"localhost:2379", "localhost:2381", "localhost:2383"},)
	serviceList,_ := cds.GetServiceList()
	fmt.Println(serviceList)
	select {}
}